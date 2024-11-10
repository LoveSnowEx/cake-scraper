package job

type Job struct {
	Company        string
	Title          string
	Link           string
	EmploymentType EmploymentType
	Seniority      Seniority
	Location       string
	NumberToHire   int
	Experience     string
	Salary         string
	Remote         Remote
	Tags           []string
	Contents       map[string]string
}

func New() *Job {
	return &Job{
		Contents: map[string]string{},
	}
}
