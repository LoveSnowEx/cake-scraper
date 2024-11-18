package app

import "cake-scraper/pkg/job"

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

func NewJob(j *job.Job) *Job {
	return &Job{
		Company:          j.Company,
		Title:            j.Title,
		Link:             j.Link,
		MainCategory:     j.MainCategory,
		SubCategory:      j.SubCategory,
		EmploymentType:   j.EmploymentType.String(),
		Seniority:        j.Seniority.String(),
		Location:         j.Location,
		NumberToHire:     j.NumberToHire,
		Experience:       j.Experience,
		Salary:           j.Salary,
		Remote:           j.Remote.String(),
		InterviewProcess: j.InterviewProcess,
		JobDescription:   j.JobDescription,
		Requirements:     j.Requirements,
		Tags:             j.Tags,
	}
}
