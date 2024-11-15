package scraper

import (
	"cake-scraper/pkg/htmlparser"
	"cake-scraper/pkg/job"
	"cake-scraper/pkg/repo/jobrepo"
	"cake-scraper/pkg/util"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly/v2"

	_ "cake-scraper/pkg/logger"
)

const (
	BackendDeveloper  Profession = "it_back-end-engineer"
	DataEngineer      Profession = "it_data-engineer"
	FrontendDeveloper Profession = "it_front-end-engineer"
	maxChanSize       int        = 100
	rateLimit                    = 30
)

var (
	_                 Scraper = (*scraper)(nil)
	jobListUrlRegex           = regexp.MustCompile(`^https://www.cake.me/jobs.*$`)
	jobDetailUrlRegex         = regexp.MustCompile(`^https://www.cake.me/companies/(.*)/jobs/(.*)$`)
	pagePattern               = regexp.MustCompile(`https://www.cake.me/jobs\?.*?page=(\d+)`)

	Cnt atomic.Int64
)

type Profession string

func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		// colly.URLFilters(jobDetailUrlRegex, jobListUrlRegex),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "locale=en")
	})
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Millisecond * 200,
		Parallelism: rateLimit,
	})
	return c
}

func (p Profession) String() string {
	return string(p)
}

func buildJobListUrl(profession Profession, page int) string {
	return fmt.Sprintf("https://www.cake.me/jobs?location_list%%5B0%%5D=Taiwan&profession%%5B0%%5D=%s&order=latest&page=%d", profession, page)
}

func parsePageNumber(url string) int {
	matches := pagePattern.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 1
	}
	page, err := strconv.Atoi(matches[1])
	if err != nil {
		util.PanicError(err)
	}
	return page
}

type Scraper interface {
	Query(conditions map[string]interface{}) []*job.Job
	Update() error
}

type scraper struct {
	Professions     []Profession
	MaxPage         int
	linkCollector   *colly.Collector
	detailCollector *colly.Collector
	jobRepo         jobrepo.JobRepo
	logger          *slog.Logger
}

func NewScraper(MaxPage int, Professions ...Profession) *scraper {
	s := &scraper{
		MaxPage:         MaxPage,
		Professions:     Professions,
		linkCollector:   NewCollector(),
		detailCollector: NewCollector(),
		jobRepo:         jobrepo.NewJobRepo(),
	}
	s.Init()
	return s
}

func (s *scraper) Init() {
	s.logger = slog.Default().WithGroup("scraper")
	s.linkCollector.OnHTML("div[class^='JobSearchHits_list__']", func(e *colly.HTMLElement) {
		hrefs := e.ChildAttrs("a[class^='JobSearchItem_jobTitle__']", "href")
		hrefs = util.Filter(hrefs, func(href string) bool {
			return href != ""
		})
		links := util.Map(hrefs, func(href string) string {
			return e.Request.AbsoluteURL(href)
		})
		for _, link := range links {
			s.handleScrapedLink(link)
		}
	})
	s.detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		Cnt.Add(1)
		j := job.New()
		j.Company = e.ChildText("div[class^='JobDescriptionLeftColumn_companyInfo__'] > a > h2")
		j.Title = e.ChildText("h1[class^='JobDescriptionLeftColumn_title__']")
		j.Link = e.Request.URL.String()
		j.Remote = job.NoRemote
		// Job Info
		e.ForEach("div[class^='JobDescriptionRightColumn_jobInfo__'] > div[class^='JobDescriptionRightColumn_row__']", func(_ int, row *colly.HTMLElement) {
			icons := util.Filter(strings.Split(row.ChildAttr("i", "class"), " "), func(str string) bool {
				return str != ""
			})
			anchors := row.ChildTexts("a")
			spans := row.ChildTexts("span")
			if len(icons) == 0 {
				// EmploymentType, Seniority, Tags
				for _, anchor := range anchors {
					if employmentType := job.NewEmploymentType(anchor); employmentType != job.InvalidEmploymentType {
						j.EmploymentType = employmentType
					} else if seniority := job.NewSeniority(anchor); seniority != job.InvalidSeniority {
						j.Seniority = seniority
					} else {
						j.Tags = append(j.Tags, anchor)
					}
				}
			} else {
				// Location, NumberToHire, Experience, Salary, Remote, Tags
				for _, icon := range icons {
					switch icon {
					case "fa-map-marker-alt":
						j.Location = anchors[0]
					case "fa-user":
						j.NumberToHire, _ = strconv.Atoi(spans[0])
					case "fa-business-time":
						j.Experience = spans[0]
					case "fa-dollar-sign":
						j.Salary = spans[0]
					case "fa-house":
						j.Remote = job.NewRemote(spans[0])
					case "fa-ellipsis-h":
						j.Tags = append(j.Tags, anchors[0])
					}
				}
			}
		})
		// Job Content
		e.ForEach("div[class^='ContentSection_contentSection__']", func(_ int, section *colly.HTMLElement) {
			contentType := section.ChildText("h3[class^='ContentSection_title__']")
			content, _ := section.DOM.Find("div[class^='RailsHtml_container__']").Html()
			content = htmlparser.Parse(content)
			if content == "" {
				return
			}
			switch contentType {
			case "Interview process":
				j.InterviewProcess = content
			case "Job Description":
				j.JobDescription = content
			case "Requirements":
				j.Requirements = content
			}
		})
		s.handleScrapedJob(j)
	})
	s.linkCollector.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 404 {
			return
		}
		s.logger.Error("linkCollector on err:", "URL", r.Request.URL, "Code", r.StatusCode, "Error", err)
	})
	s.detailCollector.OnError(func(r *colly.Response, err error) {
		s.logger.Error("detailCollector on err:", "URL", r.Request.URL, "Code", r.StatusCode, "Error", err)
	})
}

func (s *scraper) handleScrapedLink(link string) {
	if err := s.detailCollector.Visit(link); err != nil {
		util.PanicError(err)
	}
}

func (s *scraper) handleScrapedJob(j *job.Job) {
	if err := s.jobRepo.Save(j); err != nil {
		util.PanicError(err)
	}
}

func (s *scraper) Query(conditions map[string]interface{}) []*job.Job {
	jobs, err := s.jobRepo.Find(conditions)
	if err != nil {
		util.PanicError(err)
	}
	return jobs
}

func (s *scraper) Update() error {
	for _, profession := range s.Professions {
		for page := 1; page <= s.MaxPage; page++ {
			if err := s.linkCollector.Visit(buildJobListUrl(profession, page)); err != nil {
				return err
			}
		}
	}
	s.linkCollector.Wait()
	s.detailCollector.Wait()
	return nil
}
