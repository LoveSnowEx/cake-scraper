package job_test

import (
	"cake-scraper/pkg/job"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type JobTestSuite struct {
	suite.Suite
}

func (s *JobTestSuite) TestMarshalJSONAndUnmarshalJSON() {
	j := &job.Job{
		Company:        "Google",
		Title:          "Software Engineer",
		Link:           "https://google.com",
		EmploymentType: job.FullTime,
		Seniority:      job.MidSeniorLevel,
		Location:       "Mountain View, CA",
		NumberToHire:   10,
		Experience:     "Mid-Level",
		Salary:         "$100,000",
		Remote:         job.FullRemote,
		Tags:           []string{"Go", "Python"},
		Contents: map[string]string{
			"Job Description": "This is a job description",
			"Requirements":    "This is a job requirement",
		},
	}

	data, err := json.Marshal(j)
	if !s.NoError(err) {
		return
	}

	s.T().Logf("JSON: %s", string(data))

	j2 := &job.Job{}
	if !s.NoError(json.Unmarshal(data, j2)) {
		return
	}

	if !s.Equal(j, j2) {
		return
	}
}

func TestJobTestSuite(t *testing.T) {
	suite.Run(t, new(JobTestSuite))
}
