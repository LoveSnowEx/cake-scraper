package job

import "encoding/json"

type EmploymentType int
type Seniority int
type Remote int

const (
	FullTime EmploymentType = iota
	PartTime
	Internship
	Contract
	Temporary
	Volunteer
	Freelance
	InvalidEmploymentType EmploymentType = -1
)

const (
	EntryLevel Seniority = iota
	MidSeniorLevel
	Intern
	Assistant
	Director
	Executive
	InvalidSeniority Seniority = -1
)

const (
	FullRemote Remote = iota
	PartialRemote
	OptionalRemote
	NoRemote
	InvalidRemote Remote = -1
)

func NewEmploymentType(s string) EmploymentType {
	switch s {
	case "Full-time":
		return FullTime
	case "Part-time":
		return PartTime
	case "Internship":
		return Internship
	case "Contract":
		return Contract
	case "Temporary":
		return Temporary
	case "Volunteer":
		return Volunteer
	case "Freelance":
		return Freelance
	default:
		return InvalidEmploymentType
	}
}

func (et EmploymentType) String() string {
	switch et {
	case FullTime:
		return "Full-time"
	case PartTime:
		return "Part-time"
	case Internship:
		return "Internship"
	case Contract:
		return "Contract"
	case Temporary:
		return "Temporary"
	case Volunteer:
		return "Volunteer"
	case Freelance:
		return "Freelance"
	default:
		return "Invalid"
	}
}

func (et EmploymentType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + et.String() + `"`), nil
}

func (et *EmploymentType) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*et = NewEmploymentType(str)
	return nil
}

func NewSeniority(s string) Seniority {
	switch s {
	case "Entry level":
		return EntryLevel
	case "Mid-Senior level":
		return MidSeniorLevel
	case "Intern":
		return Intern
	case "Assistant":
		return Assistant
	case "Director":
		return Director
	case "Executive (VP, GM, C-Level)":
		return Executive
	default:
		return InvalidSeniority
	}
}

func (s Seniority) String() string {
	switch s {
	case EntryLevel:
		return "Entry level"
	case MidSeniorLevel:
		return "Mid-Senior level"
	case Intern:
		return "Intern"
	case Assistant:
		return "Assistant"
	case Director:
		return "Director"
	case Executive:
		return "Executive (VP, GM, C-Level)"
	default:
		return "Invalid"
	}
}

func (s *Seniority) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *Seniority) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*s = NewSeniority(str)
	return nil
}

func NewRemote(s string) Remote {
	switch s {
	case "100% Remote Work":
		return FullRemote
	case "Partial Remote Work":
		return PartialRemote
	case "Optional Remote Work":
		return OptionalRemote
	case "No Remote Work":
		return NoRemote
	default:
		return InvalidRemote
	}
}

func (r Remote) String() string {
	switch r {
	case FullRemote:
		return "100% Remote Work"
	case PartialRemote:
		return "Partial Remote Work"
	case OptionalRemote:
		return "Optional Remote Work"
	case NoRemote:
		return "No Remote Work"
	default:
		return "Invalid"
	}
}

func (r Remote) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

func (r *Remote) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*r = NewRemote(str)
	return nil
}
