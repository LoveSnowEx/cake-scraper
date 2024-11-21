package app

import (
	"cake-scraper/pkg/dto"
	"cake-scraper/pkg/job"
	"cake-scraper/pkg/repo/jobrepo"
	"cake-scraper/pkg/util"
	"cake-scraper/view"

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
	app.Use("/*", static.New("./assets"))

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
