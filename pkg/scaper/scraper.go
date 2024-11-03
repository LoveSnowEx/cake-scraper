package scaper

import (
	"cake-scraper/pkg/job"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var (
	jobListUrlRegex   = regexp.MustCompile(`^https://www.cake.me/jobs(\?.*)?$`)
	jobDetailUrlRegex = regexp.MustCompile(`^https://www.cake.me/companies/(.*)/jobs/(.*)$`)
	collectorPool     = sync.Pool{
		New: func() interface{} {
			return NewCollector()
		},
	}
)

func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("www.cake.me"),
		colly.Async(true),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "locale=en")
	})
	return c
}

type Scraper interface {
	Scrape(url string) ([]*job.Job, error)
}

type scraper struct {
}

func isJobListUrl(u string) bool {
	return jobListUrlRegex.MatchString(u)
}

func isJobDetailUrl(u string) bool {
	return jobDetailUrlRegex.MatchString(u)
}

func NewScraper() Scraper {
	s := &scraper{}
	return s
}

func (s *scraper) Scrape(url string) ([]*job.Job, error) {
	return s.scrapeJobList(url)
}

func (s *scraper) scrapeJobList(url string) ([]*job.Job, error) {
	if !isJobListUrl(url) {
		return nil, fmt.Errorf("invalid job list url: %s", url)
	}

	c := collectorPool.Get().(*colly.Collector)
	defer collectorPool.Put(NewCollector())

	jobsCh := make(chan *job.Job, 1000)

	wg := sync.WaitGroup{}
	c.OnHTML("a[data-algolia-event-name='click_job']", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		wg.Add(1)
		go func() {
			defer wg.Done()
			job, err := s.scrapeJobDetails(e.Request.AbsoluteURL(link))
			if err != nil {
				fmt.Printf("failed to scrape job details: %v\n", err)
				return
			}
			jobsCh <- job
		}()
	})
	if err := c.Visit(url); err != nil {
		return nil, err
	}
	c.Wait()
	wg.Wait()
	close(jobsCh)

	jobs := []*job.Job{}
	for job := range jobsCh {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s *scraper) scrapeJobDetails(url string) (*job.Job, error) {
	if !isJobDetailUrl(url) {
		return nil, fmt.Errorf("invalid job detail url: %s", url)
	}

	c := collectorPool.Get().(*colly.Collector)
	defer collectorPool.Put(NewCollector())

	j := job.New()
	j.Link = url
	c.OnHTML("a[class^='JobDescriptionLeftColumn_name__']", func(e *colly.HTMLElement) {
		j.Company = e.ChildText("h2")
		fmt.Printf("Company: %s\n", j.Company)
	})
	c.OnHTML("h1[class^='JobDescriptionLeftColumn_title__']", func(e *colly.HTMLElement) {
		j.Title = e.Text
		fmt.Printf("Title: %s\n", j.Title)
	})
	c.OnHTML("div[class^='JobDescriptionRightColumn_jobInfo__'] > div[class^='JobDescriptionRightColumn_row__']", func(e *colly.HTMLElement) {
		var icons, hrefs, spans []string
		e.ForEach("i", func(_ int, icon *colly.HTMLElement) {
			classes := strings.Split(icon.Attr("class"), " ")
			for _, class := range classes {
				if class == "" {
					continue
				}
				icons = append(icons, class)
			}
		})
		e.ForEach("a", func(_ int, href *colly.HTMLElement) {
			text := strings.Trim(href.Text, " ")
			if text == "" {
				return
			}
			hrefs = append(hrefs, text)
		})
		e.ForEach("span", func(_ int, span *colly.HTMLElement) {
			text := strings.Trim(span.Text, " ")
			if text == "" {
				return
			}
			spans = append(spans, text)
		})
		if len(icons) == 0 {
			// EmploymentType, Seniority, Tags
			for _, href := range hrefs {
				if employmentType := job.NewEmploymentType(href); employmentType != job.InvalidEmploymentType {
					j.Info.EmploymentType = employmentType
				} else if seniority := job.NewSeniority(href); seniority != job.InvalidSeniority {
					j.Info.Seniority = seniority
				} else {
					j.Info.Tags = append(j.Info.Tags, href)
				}
			}
		} else {
			for _, icon := range icons {
				switch icon {
				case "fa-map-marker-alt":
					j.Info.Location = hrefs[0]
				case "fa-user":
					j.Info.NumberToHire, _ = strconv.Atoi(spans[0])
				case "fa-business-time":
					j.Info.Experience = spans[0]
				case "fa-dollar-sign":
					j.Info.Salary = spans[0]
				case "fa-house":
					if remote := job.NewRemote(spans[0]); remote != job.InvalidRemote {
						j.Info.Remote = remote
					}
				case "fa-ellipsis-h":
					j.Info.Tags = append(j.Info.Tags, hrefs[0])
				}
			}
		}
	})
	c.OnHTML("div[class^='ContentSection_contentSection__']", func(e *colly.HTMLElement) {
		contentType := e.ChildText("h3[class^='ContentSection_title__']")
		content, _ := e.DOM.Find("div[class^='RailsHtml_container__']").Html()
		j.Contents[contentType] = content
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}
	c.Wait()

	return j, nil
}
