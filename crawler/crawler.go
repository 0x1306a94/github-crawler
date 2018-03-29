package crawler

import (
	"errors"
	"time"
)

type Crawler struct {
	languages 	[]string
	pool 		*workPool
	task 		[]task
	ticker  	*time.Ticker
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) Result() <-chan WorkPoolResult {
	if c.pool == nil {
		return nil
	}
	return c.pool.resultQueue
}

func (c *Crawler) AllLanguages() []string  {
	return c.languages
}

func (c *Crawler) Start() (error) {
	languages, err := c.fetchLanguage()
	if err != nil {
		return err
	}
	c.languages = []string{AllLanguage}
	c.languages = append(c.languages, languages...)
	c.pool = newWorkPool(10)
	c.task = make([]task, 0)

	for _, since := range []string{SinceToDay, SinceWeek, SinceMonth} {
		for _, lan := range languages {
			reposT := &repoTask{}
			reposT.language = lan
			reposT.since = since
			reposT.retry = 3


			developerT := &developerTask{}
			developerT.language = lan
			developerT.since = since
			developerT.retry = 3

			c.task = append(c.task, reposT, developerT)
		}
	}

	c.ticker = time.NewTicker(time.Minute * 5)
	go c.pool.run()
	go func(crawler *Crawler) {
		startOffset := 0
		sectionSize := 10
		for range crawler.ticker.C {
			var tasks []task = nil
			if startOffset + sectionSize >= len(c.task) {
				tasks = crawler.task[startOffset:]
				startOffset = 0
			} else {
				tasks = crawler.task[startOffset:(startOffset + sectionSize)]
				startOffset = startOffset + sectionSize
			}
			for _, t := range tasks {
				crawler.pool.addTask(t)
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
		if v, ok := r.([]string); ok {
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