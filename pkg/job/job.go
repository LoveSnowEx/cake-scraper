package job

import (
	"encoding/json"
)

type Job struct {
	Company  string
	Title    string
	Link     string
	Info     Info
	Contents Content
}

type Info struct {
	EmploymentType EmploymentType
	Seniority      Seniority
	Location       string
	NumberToHire   int
	Experience     string
	Salary         string
	Remote         Remote
	Tags           []string
}

type Content map[string]string

func New() *Job {
	return &Job{
		Info: Info{
			Tags: []string{},
		},
		Contents: map[string]string{},
	}
}

func (info *Info) MarshalJSON() ([]byte, error) {
	return json.Marshal(&map[string]interface{}{
		"EmploymentType": info.EmploymentType.String(),
		"Seniority":      info.Seniority.String(),
		"Location":       info.Location,
		"NumberToHire":   info.NumberToHire,
		"Experience":     info.Experience,
		"Salary":         info.Salary,
		"Remote":         info.Remote.String(),
		"Tags":           info.Tags,
	})
}
