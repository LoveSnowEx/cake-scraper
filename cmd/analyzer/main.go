package main

import (
	"cake-scraper/pkg/deeplx"
	"cake-scraper/pkg/htmlparser"
	"cake-scraper/pkg/job"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

func preprocess(jobs []*job.Job) {
	for _, j := range jobs[5:6] {
		var err error
		jd := htmlparser.Parse(j.Contents["Job Description"])
		reqs := htmlparser.Parse(j.Contents["Requirements"])
		jd, err = deeplx.Translate(jd, "ZH-TW", "EN")
		if err != nil {
			log.Fatalf("Failed to translate: %v", err)
		}
		reqs, err = deeplx.Translate(reqs, "ZH-TW", "EN")
		if err != nil {
			log.Fatalf("Failed to translate: %v", err)
		}
		j.Contents["Job Description"] = jd
		j.Contents["Requirements"] = reqs
	}
}

func main() {
	jobs, _ := readJobs()
	preprocess(jobs)
	analyzer := job.NewJobAnalyzer()
	summaries := make([]string, len(jobs))
	for _, j := range jobs {
		summary, err := analyzer.Analyze(j)
		if err != nil {
			log.Fatalf("Failed to analyze job: %v", err)
		}
		fmt.Println(summary)
		summaries = append(summaries, summary)
	}
	jsonData := "[\n" + strings.Join(summaries, ",\n") + "\n]"
	if err := os.WriteFile("summaries.json", []byte(jsonData), 0644); err != nil {
		log.Fatalf("Failed to write jobs: %v", err)
	}
}
