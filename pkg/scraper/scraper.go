package scraper

import (
	"cake-scraper/pkg/htmlparser"
	"cake-scraper/pkg/job"
	"cake-scraper/pkg/util"
	"regexp"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
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
	db        *sqlx.DB
}

func NewScraper() *scraper {
	s := &scraper{
		collector: NewCollector(),
		urls:      []string{},
	}
	return s
}

func (s *scraper) Init() {
	// Create job table
	s.db.MustExec("DROP TABLE IF EXISTS jobs;")
	s.db.MustExec(`
			CREATE TABLE jobs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				company_id TEXT NOT NULL,
				title_id TEXT NOT NULL,
				company TEXT NOT NULL DEFAULT '',
				title TEXT NOT NULL DEFAULT '',
				link TEXT NOT NULL DEFAULT '',
				employment_type INTEGER NOT NULL DEFAULT -1,
				seniority INTEGER NOT NULL DEFAULT -1,
				location TEXT NOT NULL DEFAULT '',
				number_to_hire INTEGER NOT NULL DEFAULT 0,
				experience TEXT NOT NULL DEFAULT '',
				salary TEXT NOT NULL DEFAULT '',
				remote INTEGER NOT NULL DEFAULT -1
			);
		`)
	s.db.MustExec("CREATE UNIQUE INDEX uq_jobs_company_id_title_id ON jobs (company_id, title_id);")
	// Create job_tags table
	s.db.MustExec("DROP TABLE IF EXISTS job_tags;")
	s.db.MustExec(`
		CREATE TABLE job_tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			tag TEXT NOT NULL DEFAULT ''
			CONSTRAINT fk_job_id REFERENCES jobs (id)
		);
	`)
	s.db.MustExec("CREATE UNIQUE INDEX uq_job_tags_job_id_tag ON job_tags (job_id, tag);")
	// Create job_contents table
	s.db.MustExec("DROP TABLE IF EXISTS job_contents")
	s.db.MustExec(`
		CREATE TABLE job_contents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT ''
			CONSTRAINT fk_job_id REFERENCES jobs (id)
		);
	`)
	s.db.MustExec("CREATE UNIQUE INDEX uq_job_contents_job_id_type ON job_contents (job_id, type);")

	// Scrape job list
	s.collector.OnHTML("a[class^='JobSearchItem_jobTitle__']", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		companyID, titleID := parseJobDetailUrl(link)
		sql, args := sq.Insert("jobs").
			Columns("company_id", "title_id", "link").
			Values(companyID, titleID, link).
			Suffix("ON CONFLICT DO NOTHING").
			MustSql()
		s.db.MustExec(sql, args...)
		_ = s.collector.Visit(link)
	})

	// Scrape company name
	s.collector.OnHTML("div[class^='JobDescriptionLeftColumn_companyInfo__']", func(e *colly.HTMLElement) {
		companyID, titleID := parseJobDetailUrl(e.Request.URL.String())
		comapny := e.ChildText("h2")
		sql, args := sq.Update("jobs").
			Where(sq.Eq{"company_id": companyID, "title_id": titleID}).
			Set("company", comapny).
			MustSql()
		s.db.MustExec(sql, args...)
	})

	// Scrape job title
	s.collector.OnHTML("h1[class^='JobDescriptionLeftColumn_title__']", func(e *colly.HTMLElement) {
		companyID, titleID := parseJobDetailUrl(e.Request.URL.String())
		title := e.Text
		sql, args := sq.Update("jobs").
			Where(sq.Eq{"company_id": companyID, "title_id": titleID}).
			Set("title", title).
			MustSql()
		s.db.MustExec(sql, args...)
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
		jobID := func() int64 {
			sql, args := sq.Select("id").
				From("jobs").
				Where(sq.Eq{"company_id": companyID, "title_id": titleID}).
				MustSql()
			row := s.db.QueryRowx(sql, args...)
			var id int64
			if err := row.Scan(&id); err != nil {
				panic(err)
			}
			return id
		}()
		if len(icons) == 0 {
			// EmploymentType, Seniority, Tags
			for _, anchor := range anchors {
				if employmentType := job.NewEmploymentType(anchor); employmentType != job.InvalidEmploymentType {
					sql, args := sq.Update("jobs").
						Where(sq.Eq{"id": jobID}).
						Set("employment_type", employmentType).
						MustSql()
					s.db.MustExec(sql, args...)
				} else if seniority := job.NewSeniority(anchor); seniority != job.InvalidSeniority {
					sql, args := sq.Update("jobs").
						Where(sq.Eq{"id": jobID}).
						Set("seniority", seniority).
						MustSql()
					s.db.MustExec(sql, args...)
				} else {
					tag := anchor
					sql, args := sq.Insert("job_tags").
						Columns("job_id", "tag").
						Values(jobID, tag).
						Suffix("ON CONFLICT DO NOTHING").
						MustSql()
					s.db.MustExec(sql, args...)
				}
			}
		} else {
			for _, icon := range icons {
				switch icon {
				case "fa-map-marker-alt":
					location := anchors[0]
					sql, args := sq.Update("jobs").
						Where(sq.Eq{"id": jobID}).
						Set("location", location).
						MustSql()
					s.db.MustExec(sql, args...)
				case "fa-user":
					numberToHire, _ := strconv.Atoi(spans[0])
					sql, args := sq.Update("jobs").
						Where(sq.Eq{"id": jobID}).
						Set("number_to_hire", numberToHire).
						MustSql()
					s.db.MustExec(sql, args...)
				case "fa-business-time":
					experience := spans[0]
					sql, args := sq.Update("jobs").
						Where(sq.Eq{"id": jobID}).
						Set("experience", experience).
						MustSql()
					s.db.MustExec(sql, args...)
				case "fa-dollar-sign":
					salary := spans[0]
					sql, args := sq.Update("jobs").
						Where(sq.Eq{"id": jobID}).
						Set("salary", salary).
						MustSql()
					s.db.MustExec(sql, args...)
				case "fa-house":
					if spans[0] == "" {
						remote := job.NoRemote
						sql, args := sq.Update("jobs").
							Where(sq.Eq{"id": jobID}).
							Set("remote", remote).
							MustSql()
						s.db.MustExec(sql, args...)
					} else if remote := job.NewRemote(spans[0]); remote != job.InvalidRemote {
						sql, args := sq.Update("jobs").
							Where(sq.Eq{"id": jobID}).
							Set("remote", remote).
							MustSql()
						s.db.MustExec(sql, args...)
					}
				case "fa-ellipsis-h":
					tags := anchors[0]
					sql, args := sq.Insert("job_tags").
						Columns("job_id", "tag").
						Values(jobID, tags).
						Suffix("ON CONFLICT DO NOTHING").
						MustSql()
					s.db.MustExec(sql, args...)
				}
			}
		}
	})

	// Scrape job contents
	s.collector.OnHTML("div[class^='ContentSection_contentSection__']", func(e *colly.HTMLElement) {
		contentType := e.ChildText("h3[class^='ContentSection_title__']")
		content, _ := e.DOM.Find("div[class^='RailsHtml_container__']").Html()
		companyID, titleID := parseJobDetailUrl(e.Request.URL.String())
		jobID := func() int64 {
			sql, args := sq.Select("id").
				From("jobs").
				Where(sq.Eq{"company_id": companyID, "title_id": titleID}).
				MustSql()
			row := s.db.QueryRowx(sql, args...)
			var id int64
			if err := row.Scan(&id); err != nil {
				panic(err)
			}
			return id
		}()
		sql, args := sq.Insert("job_contents").
			Columns("job_id", "type", "content").
			Values(
				jobID,
				contentType,
				htmlparser.Parse(content)).
			Suffix("ON CONFLICT (job_id, type) DO UPDATE SET content = EXCLUDED.content").
			MustSql()
		s.db.MustExec(sql, args...)
	})
}

// Query all jobs
func (s *scraper) queryJobs() ([]*job.Job, error) {
	sql, args := sq.Select("*").
		From("jobs").
		OrderBy("id").
		MustSql()
	jobRows, err := s.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	defer jobRows.Close()
	var jobs []*job.Job
	for jobRows.Next() {
		m := &struct {
			ID             int64  `db:"id"`
			CompanyID      string `db:"company_id"`
			TitleID        string `db:"title_id"`
			Company        string `db:"company"`
			Title          string `db:"title"`
			Link           string `db:"link"`
			EmploymentType int64  `db:"employment_type"`
			Seniority      int64  `db:"seniority"`
			Location       string `db:"location"`
			NumberToHire   int64  `db:"number_to_hire"`
			Experience     string `db:"experience"`
			Salary         string `db:"salary"`
			Remote         int64  `db:"remote"`
		}{}
		err := jobRows.StructScan(m)
		if err != nil {
			return nil, err
		}
		j := &job.Job{
			Company:        m.Company,
			Title:          m.Title,
			Link:           m.Link,
			EmploymentType: job.EmploymentType(m.EmploymentType),
			Seniority:      job.Seniority(m.Seniority),
			Location:       m.Location,
			NumberToHire:   int(m.NumberToHire),
			Experience:     m.Experience,
			Salary:         m.Salary,
			Remote:         job.Remote(m.Remote),
			Tags:           []string{},
			Contents:       map[string]string{},
		}
		// Queries job tags
		if err := func() error {
			sql, args := sq.Select("tag").
				From("job_tags").
				Where(sq.Eq{"job_id": m.ID}).
				MustSql()
			tagRows, err := s.db.Queryx(sql, args...)
			if err != nil {
				return err
			}
			defer tagRows.Close()
			for tagRows.Next() {
				t := map[string]interface{}{}
				err := tagRows.MapScan(t)
				if err != nil {
					return err
				}
				j.Tags = append(j.Tags, t["tag"].(string))
			}
			return nil
		}(); err != nil {
			return nil, err
		}
		// Queries job contents
		if err := func() error {
			sql, args := sq.Select("type", "content").
				From("job_contents").
				Where(sq.Eq{"job_id": m.ID}).
				MustSql()
			contentRows, err := s.db.Queryx(sql, args...)
			if err != nil {
				return err
			}
			defer contentRows.Close()
			for contentRows.Next() {
				c := map[string]interface{}{}
				err := contentRows.MapScan(c)
				if err != nil {
					return err
				}
				j.Contents[c["type"].(string)] = c["content"].(string)
			}
			return nil
		}(); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (s *scraper) AddUrl(url string) {
	s.urls = append(s.urls, url)
}

func (s *scraper) Run() []*job.Job {
	s.db = sqlx.MustConnect(sqliteshim.ShimName, "file::memory:?cache=shared")
	defer s.db.Close()
	s.Init()
	for _, url := range s.urls {
		_ = s.collector.Visit(url)
	}
	s.collector.Wait()
	jobs, err := s.queryJobs()
	if err != nil {
		panic(err)
	}
	return jobs
}
