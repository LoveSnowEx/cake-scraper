package app

import (
	"cake-scraper/pkg/dto"
	"cake-scraper/pkg/job"
	"cake-scraper/pkg/repo/jobrepo"
	"cake-scraper/pkg/util"
	"cake-scraper/view"
	jobcomponent "cake-scraper/view/components/jobs"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type App struct {
	*fiber.App
	jobRepo jobrepo.JobRepo
}

func New(app *fiber.App) *App {
	a := &App{
		app,
		jobrepo.NewJobRepo(),
	}

	app.Get("/", adaptor.HTTPHandler(
		templ.Handler(view.Index()),
	))
	app.Use("/assets/*", static.New("./assets"))
	app.Get("/components/jobs", a.JobsComponent)

	api := app.Group("/api")
	api.Get("/jobs", a.Jobs)

	return a
}

func (a *App) Jobs(c fiber.Ctx) error {
	jobs, err := a.jobRepo.Find(nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	jobsDTO := util.Map(
		jobs,
		func(j *job.Job) *dto.Job {
			return parseJob(j)
		},
	)
	return c.JSON(fiber.Map{
		"jobs": jobsDTO,
	})
}

func (a *App) JobsComponent(c fiber.Ctx) error {
	queries := c.Queries()
	conditions := jobrepo.NewConditions()
	if compony, ok := queries["company"]; ok {
		conditions.Company(compony)
	}
	if title, ok := queries["title"]; ok {
		conditions.Title(title)
	}
	if employmentTypes, ok := queries["employmentTypes"]; ok {
		for _, employmentType := range strings.Split(employmentTypes, ",") {
			if et := job.NewEmploymentType(employmentType); et != job.InvalidEmploymentType {
				conditions.EmploymentType(et)
			}
		}
	}
	if seniorities, ok := queries["seniorities"]; ok {
		for _, seniority := range strings.Split(seniorities, ",") {
			if s := job.NewSeniority(seniority); s != job.InvalidSeniority {
				conditions.Seniority(s)
			}
		}
	}
	if remotes, ok := queries["remotes"]; ok {
		for _, remote := range strings.Split(remotes, ",") {
			if r := job.NewRemote(remote); r != job.InvalidRemote {
				conditions.Remote(r)
			}
		}
	}
	if tags, ok := queries["tags"]; ok {
		conditions.Tags(tags)
	}
	paginatior := a.jobRepo.FindPaginated(conditions, 1, 10)
	return jobcomponent.
		List(util.NewPaginator(func(offset, limit int64) []*dto.Job {
			jobs := paginatior.Slice(offset, limit)
			jobsDTO := util.Map(
				jobs,
				func(j *job.Job) *dto.Job {
					return parseJob(j)
				},
			)
			return jobsDTO
		}, paginatior.CurrentPage(), paginatior.PerPage(), paginatior.Total())).
		Render(c.Context(), c)
}
