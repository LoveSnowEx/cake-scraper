package job

type Job struct {
	Company          string
	Title            string
	Link             string
	MainCategory     string
	SubCategory      string
	EmploymentType   EmploymentType
	Seniority        Seniority
	Location         string
	NumberToHire     int
	Experience       string
	Salary           string
	Remote           Remote
	InterviewProcess string
	JobDescription   string
	Requirements     string
	Tags             []string
}

func New() *Job {
	return &Job{
		EmploymentType: InvalidEmploymentType,
		Seniority:      InvalidSeniority,
		Remote:         InvalidRemote,
		Tags:           []string{},
	}
}
