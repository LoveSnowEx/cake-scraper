package main

import (
	"cake-scraper/pkg/scraper"
	"encoding/json"
	"log"
	"os"
)

func main() {
	const maxPage = 100
	professions := []scraper.Profession{
		scraper.BackendDeveloper,
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
	_ = os.MkdirAll("out", 0755)
	if err := os.WriteFile("out/jobs.json", jobsJson, 0644); err != nil {
		log.Fatal(err)
	}
}
