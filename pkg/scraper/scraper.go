package scraper

import (
	"cake-scraper/pkg/job"
	"cake-scraper/pkg/jobrepo"
	"cake-scraper/pkg/util"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var (
	jobListUrlRegex   = regexp.MustCompile(`^https://www.cake.me/jobs.*$`)
	jobDetailUrlRegex = regexp.MustCompile(`^https://www.cake.me/companies/(.*)/jobs/(.*)$`)
)

// Parser job detail url to company and title
func parseJobDetailUrl(url string) (companyID string, titleID string) {
	matches := jobDetailUrlRegex.FindStringSubmatch(url)
	return matches[1], matches[2]
}

func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.URLFilters(jobDetailUrlRegex, jobListUrlRegex),
		colly.Async(true),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "locale=en")
	})
	return c
}

type Scraper interface {
	AddUrl(url string)
	Run() []*job.Job
}

type scraper struct {
	collector *colly.Collector
	urls      []string
	repo      jobrepo.JobRepo
}

func NewScraper() *scraper {
	s := &scraper{
		collector: NewCollector(),
		urls:      []string{},
	}
	return s
}

func (s *scraper) Init() {
	s.repo.Init()

	// Scrape job list
	s.collector.OnHTML("a[class^='JobSearchItem_jobTitle__']", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		companyID, titleID := parseJobDetailUrl(link)
		_, err := s.repo.RecreateJob(companyID, titleID, link)
		util.PanicError(err)
		_ = s.collector.Visit(link)
	})

	// Scrape company name
	s.collector.OnHTML("div[class^='JobDescriptionLeftColumn_companyInfo__']", func(e *colly.HTMLElement) {
		companyID, titleID := parseJobDetailUrl(e.Request.URL.String())
		comapny := e.ChildText("h2")
		err := s.repo.UpdateJob(
			map[string]interface{}{
				"company_id": companyID,
				"title_id":   titleID,
			},
			map[string]interface{}{
				"company": comapny,
			},
		)
		util.PanicError(err)
	})

	// Scrape job title
	s.collector.OnHTML("h1[class^='JobDescriptionLeftColumn_title__']", func(e *colly.HTMLElement) {
		companyID, titleID := parseJobDetailUrl(e.Request.URL.String())
		title := e.Text
		err := s.repo.UpdateJob(
			map[string]interface{}{
				"company_id": companyID,
				"title_id":   titleID,
			},
			map[string]interface{}{
				"title": title,
			},
		)
		util.PanicError(err)
	})

	// Scrape job info
	s.collector.OnHTML("div[class^='JobDescriptionRightColumn_jobInfo__'] > div[class^='JobDescriptionRightColumn_row__']", func(e *colly.HTMLElement) {
		var icons, anchors, spans []string
		e.ForEach("i", func(_ int, icon *colly.HTMLElement) {
			classes := strings.Split(icon.Attr("class"), " ")
			classes = util.Filter(classes, func(class string) bool {
				return class != ""
			})
			icons = append(icons, classes...)
		})
		e.ForEach("a", func(_ int, anchor *colly.HTMLElement) {
			anchors = append(anchors, anchor.Text)
		})
		e.ForEach("span", func(_ int, span *colly.HTMLElement) {
			spans = append(spans, span.Text)
		})
		companyID, titleID := parseJobDetailUrl(e.Request.URL.String())
		if len(icons) == 0 {
			// EmploymentType, Seniority, Tags
			for _, anchor := range anchors {
				if employmentType := job.NewEmploymentType(anchor); employmentType != job.InvalidEmploymentType {
					err := s.repo.UpdateJob(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						map[string]interface{}{
							"employment_type": employmentType,
						},
					)
					util.PanicError(err)
				} else if seniority := job.NewSeniority(anchor); seniority != job.InvalidSeniority {
					err := s.repo.UpdateJob(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						map[string]interface{}{
							"seniority": seniority,
						},
					)
					util.PanicError(err)
				} else {
					tag := anchor
					err := s.repo.AddJobTags(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						[]string{tag},
					)
					util.PanicError(err)
				}
			}
		} else {
			for _, icon := range icons {
				switch icon {
				case "fa-map-marker-alt":
					location := anchors[0]
					err := s.repo.UpdateJob(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						map[string]interface{}{
							"location": location,
						},
					)
					util.PanicError(err)
				case "fa-user":
					numberToHire, _ := strconv.Atoi(spans[0])
					err := s.repo.UpdateJob(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						map[string]interface{}{
							"number_to_hire": numberToHire,
						},
					)
					util.PanicError(err)
				case "fa-business-time":
					experience := spans[0]
					err := s.repo.UpdateJob(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						map[string]interface{}{
							"experience": experience,
						},
					)
					util.PanicError(err)
				case "fa-dollar-sign":
					salary := spans[0]
					err := s.repo.UpdateJob(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						map[string]interface{}{
							"salary": salary,
						},
					)
					util.PanicError(err)
				case "fa-house":
					if spans[0] == "" {
						remote := job.NoRemote
						err := s.repo.UpdateJob(
							map[string]interface{}{
								"company_id": companyID,
								"title_id":   titleID,
							},
							map[string]interface{}{
								"remote": remote,
							},
						)
						util.PanicError(err)
					} else if remote := job.NewRemote(spans[0]); remote != job.InvalidRemote {
						err := s.repo.UpdateJob(
							map[string]interface{}{
								"company_id": companyID,
								"title_id":   titleID,
							},
							map[string]interface{}{
								"remote": remote,
							},
						)
						util.PanicError(err)
					}
				case "fa-ellipsis-h":
					tag := anchors[0]
					err := s.repo.AddJobTags(
						map[string]interface{}{
							"company_id": companyID,
							"title_id":   titleID,
						},
						[]string{tag},
					)
					util.PanicError(err)
				}
			}
		}
	})

	// Scrape job contents
	s.collector.OnHTML("div[class^='ContentSection_contentSection__']", func(e *colly.HTMLElement) {
		contentType := e.ChildText("h3[class^='ContentSection_title__']")
		content, _ := e.DOM.Find("div[class^='RailsHtml_container__']").Html()
		companyID, titleID := parseJobDetailUrl(e.Request.URL.String())
		err := s.repo.AddJobContent(
			map[string]interface{}{
				"company_id": companyID,
				"title_id":   titleID,
			},
			map[string]string{
				contentType: content,
			},
		)
		util.PanicError(err)
	})
}

// Query all jobs
func (s *scraper) queryJobs() ([]*job.Job, error) {
	return s.repo.FindAllJobs()
}

func (s *scraper) AddUrl(url string) {
	s.urls = append(s.urls, url)
}

func (s *scraper) Run() []*job.Job {
	db := sqlx.MustConnect(sqliteshim.ShimName, "file::memory:?cache=shared")
	s.repo = jobrepo.NewJobRepo(db)
	defer db.Close()
	s.Init()
	for _, url := range s.urls {
		_ = s.collector.Visit(url)
	}
	s.collector.Wait()
	jobs, err := s.queryJobs()
	util.PanicError(err)
	return jobs
}
