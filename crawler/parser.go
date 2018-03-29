package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
	"bufio"
	"strings"
	"errors"
)

func fetchResponseFrom(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	fmt.Println("fetch url: ", url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response statusCode not equal to 200 OK")
	}
	return resp, nil
}

func parserFromAllLanguage() []string {
	resp, err := fetchResponseFrom("https://github.com/trending")
	if err != nil {
		fmt.Println("fetchResponse error: ", err)
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(resp.Body))
	if err != nil {
		fmt.Println("NewDocumentFromReader error: ", err)
		return nil
	}

	languageSelection := doc.
		Find("div[class=select-menu-list]").
		Find("a[role=menuitem]")

	if languageSelection == nil || len(languageSelection.Nodes) == 0 {
		return nil
	}
	languages := languageSelection.Map(func(i int, selection *goquery.Selection) string {
		return strings.TrimSpace(selection.Text())
	})

	if languages == nil || len(languages) == 0 {
		return nil
	}
	return languages
}

func parserFromRepos(language, since string) *TrendingRepoResult {

	url := TrendingRepoPath
	if language != AllLanguage {
		 url = url + "/" + strings.ToLower(strings.Replace(language, " ", "-", -1))
	} else {
		url = url + "/all"
	}

	if since != "" {
		url = url + "?since=" + since
	}
	resp, err := fetchResponseFrom(url)
	if err != nil {
		fmt.Println("fetchResponse error: ", err)
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(resp.Body))
	if err != nil {
		fmt.Println("NewDocumentFromReader error: ", err)
		return nil
	}
	listSelection := doc.Find("ol[class=repo-list]").First().Children()
	if listSelection == nil || len(listSelection.Nodes) == 0 {
		return nil
	}
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
	result := new(TrendingRepoResult)
	Ranking := 1
	listSelection.Each(func(i int, selection *goquery.Selection) {
		if obj := mapObjectFunc(i, selection); obj != nil {
			obj.(*TrendingRepo).Ranking = Ranking
			result.Repos = append(result.Repos, obj.(*TrendingRepo))
			Ranking++
		}
	})
	if len(result.Repos) == 0 {
		return nil
	}

	return result
}


func parserFromDevelopers(language, since string) *TrendingDeveloperResult {

	url := TrendingDevelopersPath
	if language != AllLanguage {
		url = url + "/" + strings.ToLower(strings.Replace(language, " ", "-", -1))
	} else {
		url = url + "/all"
	}

	if since != "" {
		url = url + "?since=" + since
	}
	resp, err := fetchResponseFrom(url)
	if err != nil {
		fmt.Println("fetchResponse error: ", err)
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(resp.Body))
	if err != nil {
		fmt.Println("NewDocumentFromReader error: ", err)
		return nil
	}
	listSelection := doc.Find("ol[class=list-style-none]").First().Children()
	if listSelection == nil || len(listSelection.Nodes) == 0 {
		return nil
	}
	mapObjectFunc := func(i int, selection *goquery.Selection) interface{} {
		developer := new(TrendingDeveloper)
		name := selection.Find("div[class=mx-2]").Find(".f3.text-normal").Text()
		name = strings.TrimSpace(name)
		name = strings.Replace(name, " ", "", -1)
		name = strings.Replace(name, "\n", "", -1)
		name = strings.Replace(name, "(", ",", -1)
		name = strings.Replace(name, ")", "", -1)
		names := strings.Split(name, ",")
		if len(names) == 2 {
			developer.Login = names[0]
			developer.NickName = names[1]
		} else {
			return nil
		}

		if avatar, ok := selection.Find("div[class=mx-2]").Find("img[class=rounded-1]").Attr("src"); ok {
			developer.Avatar = avatar
		}

		repoName := selection.Find("div[class=mx-2]").Find("a").Find("span[class=repo]").Text()
		description := selection.Find("div[class=mx-2]").Find("a").Find(".repo-snipit-description.css-truncate-target").Text()

		developer.RepoName = strings.TrimSpace(repoName)
		developer.RepoDescription = strings.TrimSpace(description)
		return developer
	}
	result := new(TrendingDeveloperResult)
	Ranking := 1
	listSelection.Each(func(i int, selection *goquery.Selection) {
		if obj := mapObjectFunc(i, selection); obj != nil {
			obj.(*TrendingDeveloper).Ranking = Ranking
			result.Developers = append(result.Developers, obj.(*TrendingDeveloper))
			Ranking++
		}
	})
	if len(result.Developers) == 0 {
		return nil
	}
	result.Language = language
	result.Since = since
	return result
}
