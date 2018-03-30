package crawler

import (
	"errors"
	"time"
	"fmt"
)

type Crawler struct {
	languages 	[]string
	pool 		*workPool
	ticker  	*time.Ticker
}

func NewCrawler(languages []string) *Crawler {
	return &Crawler{
		languages: languages,
	}
}

func (c *Crawler) Result() <-chan WorkResult {
	if c.pool == nil {
		return nil
	}
	return c.pool.resultQueue
}

func (c *Crawler) AllLanguages() []string  {
	return c.languages
}

func (c *Crawler) Start() (error) {
	if c.languages == nil || len(c.languages) == 0 {
		languages, err := c.fetchLanguage()
		if err != nil {
			return err
		}
		c.languages = languages
		c.languages = []string{AllLanguage}
		c.languages = append(c.languages, languages...)
	}

	c.pool = newWorkPool(10)
	fmt.Println("start crawler ....")
	c.ticker = time.NewTicker(time.Minute * 5)
	go c.pool.run()
	go func(crawler *Crawler) {
		startOffset := 0
		sectionSize := 10
		for range crawler.ticker.C {
			chunckLans := make([]string, 0)
			if startOffset + sectionSize >= len(crawler.languages) {
				chunckLans = crawler.languages[startOffset:]
				startOffset = 0
			} else {
				chunckLans = crawler.languages[startOffset:(startOffset + sectionSize)]
				startOffset = startOffset + sectionSize
			}
			for _, lan := range chunckLans {
				for _, since := range []string{SinceToDay, SinceWeek, SinceMonth} {
					reposT := &repoTask{}
					reposT.language = lan
					reposT.since = since
					reposT.retry = 3

					developerT := &developerTask{}
					developerT.language = lan
					developerT.since = since
					developerT.retry = 3

					crawler.pool.addTask(reposT)
					crawler.pool.addTask(developerT)
				}
			}
		}

	}(c)

	return nil
}
func (c *Crawler) Stop() {
	c.ticker.Stop()
	c.pool.stop()
}

func (c *Crawler) fetchLanguage() ([]string, error) {
	w := NewWork()
	go func(w *work) {
		w.start()
	}(w)

	go func(w *work) {
		w.jobQueue <- &languageTask{
			retry: 3,
		}
	}(w)

	var result []string = nil
	for r := range w.resultQueue {
		if v, ok := r.Result.([]string); ok {
			result = v
		}
		go func(w *work) {
			w.stop()
		}(w)
	}
	if result == nil || len(result) == 0 {
		return nil, errors.New("fetchLanguage error")
	}
	return result, nil
}