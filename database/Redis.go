package database

import (
	"github.com/go-redis/redis"
	"time"
	"errors"
	"fmt"
	"encoding/json"
)

type Redis struct {
	client  	*redis.Client
	Addr 		string
	Passwprd 	string
	DB 			int
}

func NewRedis(addr, password string, db int) *Redis {
	return &Redis{
		Addr: addr,
		Passwprd: password,
		DB: db,
	}
}

func (r *Redis) Connection() (err error) {
	r.client = redis.NewClient(&redis.Options{
		Addr: r.Addr,
		Password: r.Passwprd,
		DB: r.DB,
	})

	retryCount := 10
	count := 1
	fmt.Println("Redis start Ping")
	for {
		_, err = r.client.Ping().Result()
		if err != nil && count > retryCount {
			r.client.Close()
			return
		} else if err != nil {
			fmt.Println("Redis Ping error: ", err)
		} else {
			break
		}
		count++
		time.Sleep(time.Second * 5)
	}
	return
}
func (r *Redis) Disconnect()  {
	if r.client != nil {
		r.client.Close()
	}
}
func (r *Redis) Save(key string, obj string, expiration time.Duration) (err error) {

	err = r.client.Set(key, obj, expiration).Err()
	if err != nil {
		return
	}
	return
}

func (r *Redis) Get(key string) (result string, err error) {

	result, err = r.client.Get(key).Result()
	if err == redis.Nil {
		err = errors.New(fmt.Sprintf("key %v does not exists", key))
		return
	} else if err != nil {
		return
	}
	return
}

func (r *Redis) GetCacheLanguages() []string {

	result, err := r.Get("all-language")
	if err != nil {
		return nil
	}

	var lans []string
	err = json.Unmarshal([]byte(result), &lans)
	if err != nil {
		return nil
	}
	return lans
}

func (r *Redis) SaveLanguages(languages []string) {
	if len(languages) == 0 {
		return
	}
	data, err := json.Marshal(languages)
	if err != nil {
		return
	}
	if len(data) == 0 {
		return
	}
	err = r.Save("all-language", string(data), 0)
	if err == nil {
		fmt.Println("all-language save to redis success")
	} else {
		fmt.Println("all-language save to redis failure: ", err)
	}
}