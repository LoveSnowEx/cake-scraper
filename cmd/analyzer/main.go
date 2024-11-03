package main

import (
	"cake-scraper/pkg/htmlparser"
	"cake-scraper/pkg/job"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func readJobs() ([]*job.Job, error) {
	jobs := []*job.Job{}

	data, err := os.ReadFile("jobs.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func main() {
	jobs, _ := readJobs()
	summaries := map[string]string{}
	for _, j := range jobs {
		j.Contents["Job Description"] = htmlparser.Parse(j.Contents["Job Description"])
		j.Contents["Requirements"] = htmlparser.Parse(j.Contents["Requirements"])
		analyzer := job.NewJobAnalyzer()
		summary, err := analyzer.Analyze(j)
		if err != nil {
			log.Fatalf("Failed to analyze job: %v", err)
		}
		summaries[fmt.Sprintf("%s - %s", j.Company, j.Title)] = summary
	}
	summariesJSON, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal summaries: %v", err)
	}
	if err := os.WriteFile("summaries.json", summariesJSON, 0644); err != nil {
		log.Fatalf("Failed to write summaries: %v", err)
	}
}
