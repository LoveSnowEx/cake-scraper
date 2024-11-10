package main

import (
	sc "cake-scraper/pkg/scraper"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const jobListUrl = "https://www.cake.me/jobs?location_list%5B0%5D=Taipei%20City%2C%20Taiwan&profession%5B0%5D=it_back-end-engineer&profession%5B1%5D=it_data-engineer&job_type%5B0%5D=full_time&year_of_seniority%5B0%5D=1_3&salary_type=per_month&salary_currency=TWD&salary_range%5Bmin%5D=60000"

func main() {
	scraper := sc.NewScraper()
	scraper.AddUrl(jobListUrl)

	jobs := scraper.Run()

	for _, job := range jobs {
		fmt.Printf("Company: %s\nTitle: %s\nLink: %s\n", job.Company, job.Title, job.Link)
	}

	jobsJson, err := json.MarshalIndent(jobs, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	_ = os.MkdirAll("out", 0755)
	if err := os.WriteFile("out/jobs.json", jobsJson, 0644); err != nil {
		log.Fatal(err)
	}
}
