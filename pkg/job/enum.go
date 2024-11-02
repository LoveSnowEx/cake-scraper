package job

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

func NewRemote(s string) Remote {
	switch s {
	case "Full remote":
		return FullRemote
	case "Partial remote":
		return PartialRemote
	case "Optional remote":
		return OptionalRemote
	case "No remote":
		return NoRemote
	default:
		return InvalidRemote
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

func (r Remote) String() string {
	switch r {
	case FullRemote:
		return "Full remote"
	case PartialRemote:
		return "Partial remote"
	case OptionalRemote:
		return "Optional remote"
	case NoRemote:
		return "No remote"
	default:
		return "Invalid"
	}
}
