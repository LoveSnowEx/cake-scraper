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

func (info *Info) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	for key, value := range v {
		switch key {
		case "EmploymentType":
			info.EmploymentType = NewEmploymentType(value.(string))
		case "Seniority":
			info.Seniority = NewSeniority(value.(string))
		case "Location":
			info.Location = value.(string)
		case "NumberToHire":
			info.NumberToHire = int(value.(float64))
		case "Experience":
			info.Experience = value.(string)
		case "Salary":
			info.Salary = value.(string)
		case "Remote":
			info.Remote = NewRemote(value.(string))
		case "Tags":
			info.Tags = func() []string {
				tags := []string{}
				for _, tag := range value.([]interface{}) {
					tags = append(tags, tag.(string))
				}
				return tags
			}()
		}
	}
	return nil
}
