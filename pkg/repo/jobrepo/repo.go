package jobrepo

import (
	"cake-scraper/pkg/database"
	"cake-scraper/pkg/job"
	"cake-scraper/pkg/location"
	"cake-scraper/pkg/util"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

var _ JobRepo = (*jobRepoImpl)(nil)

type Time time.Time

type JobPo struct {
	ID               int64  `db:"id"`
	Link             string `db:"link"`
	Company          string `db:"company"`
	Title            string `db:"title"`
	EmploymentType   int64  `db:"employment_type"`
	Seniority        int64  `db:"seniority"`
	Location         string `db:"location"`
	NumberToHire     int64  `db:"number_to_hire"`
	Experience       string `db:"experience"`
	Salary           string `db:"salary"`
	Remote           int64  `db:"remote"`
	InterviewProcess string `db:"interview_process"`
	JobDescription   string `db:"job_description"`
	Requirements     string `db:"requirements"`
	CreatedAt        Time   `db:"created_at"`
	UpdatedAt        Time   `db:"updated_at"`
}

type JobRepo interface {
	Find(conditions map[string]interface{}) ([]*job.Job, error)
	FindPaginated(conditions Conditions, page, perPage int64) util.Paginator[*job.Job]
	Save(j *job.Job) error
	Delete(conditions map[string]interface{}) error
}

type jobRepoImpl struct {
	db *database.DB
}

func (t Time) Value() (time.Time, error) {
	return time.Time(t), nil
}

func (t *Time) Scan(v interface{}) error {
	if v == nil {
		*t = Time(time.Time{})
		return nil
	}
	tp, err := time.Parse(time.DateTime, v.(string))
	if err != nil {
		return err
	}
	*t = Time(tp)
	return nil
}

func (j *JobPo) ToJob() *job.Job {
	return &job.Job{
		Company:          j.Company,
		Title:            j.Title,
		Link:             j.Link,
		EmploymentType:   job.EmploymentType(j.EmploymentType),
		Seniority:        job.Seniority(j.Seniority),
		Location:         j.Location,
		NumberToHire:     int(j.NumberToHire),
		Experience:       j.Experience,
		Salary:           j.Salary,
		Remote:           job.Remote(j.Remote),
		InterviewProcess: j.InterviewProcess,
		JobDescription:   j.JobDescription,
		Requirements:     j.Requirements,
	}
}

func NewJobRepo() *jobRepoImpl {
	db, err := database.Connect()
	util.PanicError(err)
	return &jobRepoImpl{db: db}
}

func (r *jobRepoImpl) Find(conditions map[string]interface{}) ([]*job.Job, error) {
	if conditions == nil {
		conditions = map[string]interface{}{}
	}
	sql, args, err := sq.Select("*").
		From("jobs").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, err
	}
	var jobPos []*JobPo
	err = r.db.Select(&jobPos, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select jobs: %w", err)
	}
	var result []*job.Job
	for _, jobPo := range jobPos {
		j := jobPo.ToJob()
		tags, err := r.findTags(jobPo.ID)
		if err != nil {
			return nil, err
		}
		j.Tags = tags
		result = append(result, j)
	}
	return result, nil
}

func (r *jobRepoImpl) findTags(jobID int64) ([]string, error) {
	sql, args, err := sq.Select("tag").
		From("jobs_tags").
		Join("tags ON jobs_tags.tag_id = tags.id").
		Where(sq.Eq{"job_id": jobID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var tags []string
	err = r.db.Select(&tags, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	return tags, nil
}

func (r *jobRepoImpl) FindPaginated(conditions Conditions, page, perPage int64) util.Paginator[*job.Job] {
	builder := conditions.ToSelectBuilder()
	var total int64
	sql, args, err := builder.Columns("COUNT(*)").
		ToSql()
	util.PanicError(err)
	err = r.db.Get(&total, sql, args...)
	util.PanicError(err)
	return util.NewPaginator(func(offset, limit int64) []*job.Job {
		var jobPos []*JobPo
		sql, args, err := builder.Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
		util.PanicError(err)
		err = r.db.Select(&jobPos, sql, args...)
		util.PanicError(err)
		var result []*job.Job
		for _, jobPo := range jobPos {
			j := jobPo.ToJob()
			tags, err := r.findTags(jobPo.ID)
			util.PanicError(err)
			j.Tags = tags
			result = append(result, j)
		}
		return result
	}, int64(page), int64(perPage), total)
}

func (r *jobRepoImpl) Save(j *job.Job) (err error) {
	tx := r.db.MustBegin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	// Save job
	sql, args, err := sq.Insert("jobs").
		SetMap(map[string]interface{}{
			"link":              j.Link,
			"company":           j.Company,
			"title":             j.Title,
			"employment_type":   j.EmploymentType,
			"seniority":         j.Seniority,
			"location":          j.Location,
			"number_to_hire":    j.NumberToHire,
			"experience":        j.Experience,
			"salary":            j.Salary,
			"remote":            j.Remote,
			"interview_process": j.InterviewProcess,
			"job_description":   j.JobDescription,
			"requirements":      j.Requirements,
		}).
		Suffix(`
			ON CONFLICT(link) DO UPDATE SET
				title = EXCLUDED.title,
				employment_type = EXCLUDED.employment_type,
				seniority = EXCLUDED.seniority,
				location = EXCLUDED.location,
				number_to_hire = EXCLUDED.number_to_hire,
				experience = EXCLUDED.experience,
				salary = EXCLUDED.salary,
				remote = EXCLUDED.remote,
				interview_process = EXCLUDED.interview_process,
				job_description = EXCLUDED.job_description,
				requirements = EXCLUDED.requirements,
				updated_at = CURRENT_TIMESTAMP
		`).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}
	var jobID int64
	if err = tx.Get(&jobID, sql, args...); err != nil {
		return fmt.Errorf("failed to insert job: %w", err)
	}
	// Save categories
	sql, args, err = sq.Delete("jobs_categories").
		Where(sq.Eq{"job_id": jobID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete jobs_categories: %w", err)
	}
	sql, args, err = sq.Insert("categories").
		Columns("main", "sub").
		Values(j.MainCategory, j.SubCategory).
		Suffix("ON CONFLICT(main, sub) DO UPDATE SET main = EXCLUDED.main, sub = EXCLUDED.sub").
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}
	var categoryID int64
	if err = tx.Get(&categoryID, sql, args...); err != nil {
		return fmt.Errorf("failed to insert category: %w", err)
	}
	sql, args, err = sq.Insert("jobs_categories").
		Columns("job_id", "category_id").
		Values(jobID, categoryID).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to insert jobs_categories: %w", err)
	}
	// Save tags
	sql, args, err = sq.Delete("jobs_tags").
		Where(sq.Eq{"job_id": jobID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete jobs_tags: %w", err)
	}
	for _, tag := range j.Tags {
		sql, args, err := sq.Insert("tags").
			Columns("tag").
			Values(tag).
			Suffix("ON CONFLICT(tag) DO UPDATE SET tag = EXCLUDED.tag").
			Suffix("RETURNING id").
			ToSql()
		if err != nil {
			return err
		}
		var tagID int64
		if err = tx.Get(&tagID, sql, args...); err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}
		sql, args, err = sq.Insert("jobs_tags").
			Columns("job_id", "tag_id").
			Values(jobID, tagID).
			Suffix("ON CONFLICT DO NOTHING").
			ToSql()
		if err != nil {
			return err
		}
		_, err = tx.Exec(sql, args...)
		if err != nil {
			return fmt.Errorf("failed to insert jobs_tags: %w", err)
		}
	}
	// Save location
	sql, args, err = sq.Delete("jobs_locations").
		Where(sq.Eq{"job_id": jobID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete jobs_locations: %w", err)
	}
	matchedLocation := location.FindBestMatch(j.Location)
	if matchedLocation != nil {
		var locationID int64
		sql, args, err := sq.Select("id").
			From("locations").
			Where(sq.Eq{"address": matchedLocation.Address()}).
			ToSql()
		if err != nil {
			return err
		}
		if err := tx.Get(&locationID, sql, args...); err != nil {
			return fmt.Errorf("failed to select location: %w", err)
		}
		sql, args, err = sq.Insert("jobs_locations").
			Columns("job_id", "location_id").
			Values(jobID, locationID).
			Suffix("ON CONFLICT DO NOTHING").
			ToSql()
		if err != nil {
			return err
		}
		_, err = tx.Exec(sql, args...)
		if err != nil {
			return fmt.Errorf("failed to insert jobs_locations: %w", err)
		}
	}
	return nil
}

func (r *jobRepoImpl) Delete(conditions map[string]interface{}) error {
	sql, args, err := sq.Delete("jobs").
		Where(conditions).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete jobs: %w", err)
	}
	return nil
}
