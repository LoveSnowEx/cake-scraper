package dto

import "cake-scraper/pkg/util"

type Job struct {
	Company          string   `json:"company"`
	Title            string   `json:"title"`
	Link             string   `json:"link"`
	MainCategory     string   `json:"main_category"`
	SubCategory      string   `json:"sub_category"`
	EmploymentType   string   `json:"employment_type"`
	Seniority        string   `json:"seniority"`
	Location         string   `json:"location"`
	NumberToHire     int      `json:"number_to_hire"`
	Experience       string   `json:"experience"`
	Salary           string   `json:"salary"`
	Remote           string   `json:"remote"`
	InterviewProcess string   `json:"interview_process"`
	JobDescription   string   `json:"job_description"`
	Requirements     string   `json:"requirements"`
	Tags             []string `json:"tags"`
}

type JobsPaginator = util.Paginator[*Job]
