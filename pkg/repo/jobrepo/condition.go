package jobrepo

import (
	"cake-scraper/pkg/job"

	sq "github.com/Masterminds/squirrel"
)

type Conditions struct {
	company         string
	title           string
	employmentTypes []job.EmploymentType
	seniorities     []job.Seniority
	remotes         []job.Remote
	tags            []string
}

func NewConditions() Conditions {
	return Conditions{}
}

func (c Conditions) Clone() Conditions {
	return Conditions{
		company:         c.company,
		title:           c.title,
		employmentTypes: append([]job.EmploymentType{}, c.employmentTypes...),
		seniorities:     append([]job.Seniority{}, c.seniorities...),
		remotes:         append([]job.Remote{}, c.remotes...),
		tags:            append([]string{}, c.tags...),
	}
}

func (c Conditions) Company(company string) Conditions {
	clone := c.Clone()
	clone.company = company
	return clone
}

func (c Conditions) Title(title string) Conditions {
	clone := c.Clone()
	clone.title = title
	return clone
}

func (c Conditions) EmploymentType(employmentType ...job.EmploymentType) Conditions {
	clone := c.Clone()
	clone.employmentTypes = append(clone.employmentTypes, employmentType...)
	return clone
}

func (c Conditions) Seniority(senioritys ...job.Seniority) Conditions {
	clone := c.Clone()
	clone.seniorities = append(clone.seniorities, senioritys...)
	return clone
}

func (c Conditions) Remote(remote ...job.Remote) Conditions {
	clone := c.Clone()
	clone.remotes = append(clone.remotes, remote...)
	return clone
}

func (c Conditions) Tags(tags ...string) Conditions {
	c.tags = append(c.tags, tags...)
	return c
}

func (c Conditions) ToSelectBuilder(columns ...string) sq.SelectBuilder {
	builder := sq.Select(columns...).
		From("jobs AS j").
		Join("jobs_tags AS jt ON j.id = jt.job_id").
		Join("tags AS t ON jt.tag_id = t.id")

	if c.company != "" {
		builder = builder.Where(sq.Eq{"j.company": c.company})
	}
	if c.title != "" {
		builder = builder.Where(sq.Eq{"j.title": c.title})
	}
	if len(c.employmentTypes) > 0 {
		builder = builder.Where(sq.Eq{"j.employment_type": c.employmentTypes})
	}
	if len(c.seniorities) > 0 {
		builder = builder.Where(sq.Eq{"j.seniority": c.seniorities})
	}
	if len(c.remotes) > 0 {
		builder = builder.Where(sq.Eq{"j.remote": c.remotes})
	}
	if len(c.tags) > 0 {
		builder = builder.Where(sq.Eq{"t.tag": c.tags})
	}
	return builder
}
