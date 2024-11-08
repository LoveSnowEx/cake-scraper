package scraper

import "cake-scraper/pkg/job"

type result map[string]*job.Job

type resultUpdater func(*result)
