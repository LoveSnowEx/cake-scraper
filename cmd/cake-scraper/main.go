package main

import (
	"cake-scraper/pkg/scraper"
	"cake-scraper/pkg/util"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var jobJsonPath = filepath.Join(util.ProjectRoot, "out/jobs.json")

func main() {
	const maxPage = 15
	professions := []scraper.Profession{
		scraper.BackendDeveloper,
		scraper.DataEngineer,
		scraper.FrontendDeveloper,
	}
	sc := scraper.NewScraper(maxPage, professions...)
	if err := sc.Update(); err != nil {
		panic(err)
	}
	jobs := sc.Query(nil)
	jobsJson, err := json.MarshalIndent(jobs, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	_ = os.MkdirAll(filepath.Dir(jobJsonPath), 0755)
	if err := os.WriteFile(jobJsonPath, jobsJson, 0644); err != nil {
		log.Fatal(err)
	}
}
