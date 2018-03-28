package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"net/http"
	"bufio"
)

func main() {
	resp, err := http.Get("https://github.com/trending/objective-c?since=weekly")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求失败...")
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(resp.Body))
	if err != nil {
		log.Fatal(err)
	}
	listSelection := doc.Find("ol[class=repo-list]").First().Children()
	mapObjectFunc := func(i int, selection *goquery.Selection) interface{} {
		repo := new(TrendingRepo)
		fullName := selection.Find("h3").First().Text()
		description := selection.Find("div[class=py-1]").Find(".col-9.d-inline-block.text-gray.m-0.pr-4").First().Text()
		language := selection.Find("span[itemprop=programmingLanguage]").First().Text()
		stars := selection.Find(".f6.text-gray.mt-2").Find(".muted-link.d-inline-block.mr-3").First().Text()
		forkers := selection.Find(".f6.text-gray.mt-2").Find(".muted-link.d-inline-block.mr-3").Eq(1).Text()
		gains := selection.Find(".f6.text-gray.mt-2").Find(".d-inline-block.float-sm-right").First().Text()

		repo.FullName = strings.TrimSpace(strings.Replace(fullName, " ", "", -1))
		repo.Description = strings.TrimSpace(description)
		repo.Language = strings.TrimSpace(language)
		repo.Stars = strings.TrimSpace(strings.Replace(stars, ",", "", -1))
		repo.Forkers = strings.TrimSpace(strings.Replace(forkers, ",", "", -1))
		repo.Gains = strings.TrimSpace(strings.Replace(gains, ",", "", -1))
		return repo
	}
	result := make([]interface{}, 0)
	listSelection.Each(func(i int, selection *goquery.Selection) {
		result = append(result, mapObjectFunc(i, selection))
	})
	fmt.Println(result)
}
