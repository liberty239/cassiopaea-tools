package forum

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	scrapeDelay   = 3 * time.Second
	scrapeTimeout = 5 * time.Second
)

type scraper struct {
	onSessionLink func(e *colly.HTMLElement) error
	onSession     func(e *colly.HTMLElement) error
}

func (s *scraper) do(urls ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = r
			default:
				err = errors.New("panic")
			}
		}
	}()

	c := colly.NewCollector(
		colly.AllowedDomains("cassiopaea.org"),
		colly.UserAgent(
			fmt.Sprintf(
				"Mozilla/5.0 (compatible; %s)",
				packageName(),
			),
		),
		colly.Async(false),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*cassiopaea.org",
		Parallelism: 1,
	})
	c.SetRequestTimeout(scrapeTimeout)

	c.OnHTML(
		`.pageNav-jump`,
		func(e *colly.HTMLElement) {
			if isSearchPath(e.Request.URL.Path) {
				e.Request.Visit(e.Attr("href"))
			}
		})

	c.OnHTML(
		`.button--link`,
		func(e *colly.HTMLElement) {
			e.Request.Visit(e.Attr("href"))
		})

	var links sync.Map
	c.OnHTML(
		`li.block-row     > 
		 div:nth-child(1) > 
		 div:nth-child(2) > 
		 h3:nth-child(1)  > 
		 a:nth-child(1)`,
		func(e *colly.HTMLElement) {
			if _, ok := links.LoadOrStore(e.Attr("href"), nil); !ok {
				if s.onSessionLink != nil {
					if err := s.onSessionLink(e); err != nil {
						panic(err)
					}
				}
				if s.onSession != nil {
					e.Request.Visit(e.Attr("href"))
				}
			}
		})

	c.OnHTML(
		`div:nth-child(2) > 
		 div:nth-child(2) > 
		 div:nth-child(1) > 
		 div:nth-child(2) > 
		 div:nth-child(1) > 
		 article:nth-child(1) > 
		 div:nth-child(1)`,
		func(e *colly.HTMLElement) {
			if s.onSession != nil {
				if err := s.onSession(e); err != nil {
					panic(err)
				}
			}
		})

	c.OnError(func(_ *colly.Response, err error) {
		panic(err)
	})

	for _, url := range urls {
		if err = c.Visit(url); err != nil {
			return
		}
	}
	c.Wait()
	return
}
