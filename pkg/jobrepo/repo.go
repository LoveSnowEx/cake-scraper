package jobrepo

import (
	"cake-scraper/pkg/job"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var _ JobRepo = (*jobRepoImpl)(nil)

type JobPo struct {
	ID             int64  `db:"id"`
	CompanyID      string `db:"company_id"`
	TitleID        string `db:"title_id"`
	Company        string `db:"company"`
	Title          string `db:"title"`
	Link           string `db:"link"`
	Breadcrumbs    string `db:"breadcrumbs"`
	EmploymentType int64  `db:"employment_type"`
	Seniority      int64  `db:"seniority"`
	Location       string `db:"location"`
	NumberToHire   int64  `db:"number_to_hire"`
	Experience     string `db:"experience"`
	Salary         string `db:"salary"`
	Remote         int64  `db:"remote"`
}

type JobTagPo struct {
	ID    int64  `db:"id"`
	JobID int64  `db:"job_id"`
	Tag   string `db:"tag"`
}

type JobContentPo struct {
	ID      int64  `db:"id"`
	JobID   int64  `db:"job_id"`
	Type    string `db:"type"`
	Content string `db:"content"`
}

type JobRepo interface {
	Init() error
	FindAllJobs() ([]*job.Job, error)
	RecreateJob(companyID, titleID, link string) (int64, error)
	UpdateJob(conditions map[string]interface{}, values map[string]interface{}) error
	AddJobTags(condition map[string]interface{}, tags []string) error
	AddJobContent(condition map[string]interface{}, content map[string]string) error
	DeleteJob(conditions map[string]interface{}) error
}

type jobRepoImpl struct {
	db *sqlx.DB
}

func NewJobRepo(db *sqlx.DB) *jobRepoImpl {
	return &jobRepoImpl{db: db}
}

func (r *jobRepoImpl) Init() (err error) {
	// Create jobs table
	_, err = r.db.Exec("DROP TABLE IF EXISTS jobs;")
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`
		CREATE TABLE jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_id TEXT NOT NULL,
			title_id TEXT NOT NULL,
			company TEXT NOT NULL DEFAULT '',
			title TEXT NOT NULL DEFAULT '',
			link TEXT NOT NULL DEFAULT '',
			breadcrumbs TEXT NOT NULL DEFAULT '',
			employment_type INTEGER NOT NULL DEFAULT -1,
			seniority INTEGER NOT NULL DEFAULT -1,
			location TEXT NOT NULL DEFAULT '',
			number_to_hire INTEGER NOT NULL DEFAULT 0,
			experience TEXT NOT NULL DEFAULT '',
			salary TEXT NOT NULL DEFAULT '',
			remote INTEGER NOT NULL DEFAULT -1
		);
	`)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("CREATE UNIQUE INDEX uq_jobs_company_id_title_id ON jobs (company_id, title_id);")
	if err != nil {
		return err
	}
	// Create job_tags table
	_, err = r.db.Exec("DROP TABLE IF EXISTS job_tags;")
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`
		CREATE TABLE job_tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			tag TEXT NOT NULL DEFAULT ''
			CONSTRAINT fk_job_id REFERENCES jobs (id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("CREATE UNIQUE INDEX uq_job_tags_job_id_tag ON job_tags (job_id, tag);")
	if err != nil {
		return err
	}
	// Create job_contents table
	_, err = r.db.Exec("DROP TABLE IF EXISTS job_contents")
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`
		CREATE TABLE job_contents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT ''
			CONSTRAINT fk_job_id REFERENCES jobs (id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("CREATE UNIQUE INDEX uq_job_contents_job_id_type ON job_contents (job_id, type);")
	if err != nil {
		return err
	}
	return nil
}

func (r *jobRepoImpl) FindAllJobs() ([]*job.Job, error) {
	sql, args, err := sq.Select("*").From("jobs").ToSql()
	if err != nil {
		return nil, err
	}
	var jobPos []*JobPo
	err = r.db.Select(&jobPos, sql, args...)
	if err != nil {
		return nil, err
	}
	var result []*job.Job
	for _, jobPo := range jobPos {
		j := &job.Job{
			Company:        jobPo.Company,
			Title:          jobPo.Title,
			Link:           jobPo.Link,
			EmploymentType: job.EmploymentType(jobPo.EmploymentType),
			Seniority:      job.Seniority(jobPo.Seniority),
			Location:       jobPo.Location,
			NumberToHire:   int(jobPo.NumberToHire),
			Experience:     jobPo.Experience,
			Salary:         jobPo.Salary,
			Remote:         job.Remote(jobPo.Remote),
		}
		tags, err := r.findJobTags(jobPo.ID)
		if err != nil {
			return nil, err
		}
		j.Tags = tags
		contents, err := r.findJobContents(jobPo.ID)
		if err != nil {
			return nil, err
		}
		j.Contents = contents
		result = append(result, j)
	}
	return result, nil
}

func (r *jobRepoImpl) RecreateJob(companyID, titleID, link string) (int64, error) {
	err := r.DeleteJob(map[string]interface{}{"company_id": companyID, "title_id": titleID})
	if err != nil {
		return 0, err
	}
	sql, args, err := sq.Insert("jobs").
		Columns("company_id", "title_id", "link").
		Values(companyID, titleID, link).
		Suffix("ON CONFLICT(company_id, title_id) DO UPDATE SET link = EXCLUDED.link").
		ToSql()
	if err != nil {
		return 0, err
	}
	_, err = r.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	job, err := r.findJobPo(map[string]interface{}{"company_id": companyID, "title_id": titleID})
	if err != nil {
		return 0, err
	}
	return job.ID, nil
}

func (r *jobRepoImpl) UpdateJob(conditions map[string]interface{}, values map[string]interface{}) error {
	builder := sq.Update("jobs")
	for k, v := range conditions {
		builder = builder.Where(sq.Eq{k: v})
	}
	for k, v := range values {
		builder = builder.Set(k, v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(sql, args...)
	return err
}

func (r *jobRepoImpl) AddJobTags(condition map[string]interface{}, tags []string) error {
	job, err := r.findJobPo(condition)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		sql, args, err := sq.Insert("job_tags").
			Columns("job_id", "tag").
			Values(job.ID, tag).
			Suffix("ON CONFLICT(job_id, tag) DO NOTHING").
			ToSql()
		if err != nil {
			return err
		}
		_, err = r.db.Exec(sql, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *jobRepoImpl) AddJobContent(condition map[string]interface{}, content map[string]string) error {
	job, err := r.findJobPo(condition)
	if err != nil {
		return err
	}
	for k, v := range content {
		sql, args, err := sq.Insert("job_contents").
			Columns("job_id", "type", "content").
			Values(job.ID, k, v).
			Suffix("ON CONFLICT(job_id, type) DO UPDATE SET content = EXCLUDED.content").
			ToSql()
		if err != nil {
			return err
		}
		_, err = r.db.Exec(sql, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *jobRepoImpl) DeleteJob(conditions map[string]interface{}) error {
	builder := sq.Delete("jobs")
	for k, v := range conditions {
		builder = builder.Where(sq.Eq{k: v})
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(sql, args...)
	return err
}

func (r *jobRepoImpl) findJobPo(conditions map[string]interface{}) (*JobPo, error) {
	builder := sq.Select("*").From("jobs")
	for k, v := range conditions {
		builder = builder.Where(sq.Eq{k: v})
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var job JobPo
	err = r.db.Get(&job, sql, args...)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepoImpl) findJobTags(jobID int64) ([]string, error) {
	sql, args, err := sq.Select("tag").From("job_tags").Where(sq.Eq{"job_id": jobID}).ToSql()
	if err != nil {
		return nil, err
	}
	var tags []string
	err = r.db.Select(&tags, sql, args...)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *jobRepoImpl) findJobContents(jobID int64) (map[string]string, error) {
	sql, args, err := sq.Select("type", "content").
		From("job_contents").
		Where(sq.Eq{"job_id": jobID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var contentPos []*JobContentPo
	err = r.db.Select(&contentPos, sql, args...)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	for _, contentPo := range contentPos {
		result[contentPo.Type] = contentPo.Content
	}
	return result, nil
}
