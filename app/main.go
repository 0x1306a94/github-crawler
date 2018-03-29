package main

import (
	"github-crawler/crawler"
	"fmt"
	"log"
)

func main() {

	app := crawler.NewCrawler()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(app.AllLanguages())
	for result := range app.Result() {
		if result.Result != nil {
			fmt.Println(result.Result)
		}
	}
}
