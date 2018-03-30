package main

import (
	"github-crawler/crawler"
	"github-crawler/database"
	"fmt"
	"encoding/json"
	"strings"
	"time"
	"os"
)

var redisClient *database.Redis


func main() {

	redisHost := os.Getenv("RedisHost")
	redisPassword := os.Getenv("RedisPassword")
	if redisHost == "" {
		redisHost = "127.0.0.1:6379"
	}
	redisClient = database.NewRedis(redisHost, redisPassword, 0)

	if err := redisClient.Connection(); err != nil {
		fmt.Println("Redis connect error: ", err)
		return
	}

	defer redisClient.Disconnect()

	cacheLanguages := redisClient.GetCacheLanguages()
	if len(cacheLanguages) > 0 {
		fmt.Println("cacheLanguages: \n", cacheLanguages)
	}
	app := crawler.NewCrawler(cacheLanguages)
	if err := app.Start(); err != nil {
		fmt.Println("crawler start error: ", err)
		return
	}
	go redisClient.SaveLanguages(app.AllLanguages())
	for result := range app.Result() {
		if result.Result != nil && result.ResultType != crawler.ResultTypeLanguage {
			go SaveRedis(result)
		}
	}
}


func SaveRedis(result crawler.WorkResult)  {

	lan := strings.ToLower(strings.Replace(result.Language, " ", "-", -1))
	key := fmt.Sprintf("trending-%v-%v-%v", result.ResultType,lan, result.Since)
	data, err := json.Marshal(result.Result)
	if err != nil {
		return
	}
	if len(data) == 0 {
		return
	}

	err = redisClient.Save(key, string(data), time.Hour * 1)
	if err == nil {
		fmt.Printf("%v save to redis success\n", key)
	} else {
		fmt.Printf("%v save to redis failure: %v\n", key, err.Error())
	}
}

