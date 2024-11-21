package app

import (
	"cake-scraper/pkg/dto"
	"cake-scraper/pkg/job"
)

func parseJob(j *job.Job) *dto.Job {
	return &dto.Job{
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
