package locationrepo

import (
	"cake-scraper/pkg/database"
	"cake-scraper/pkg/location"

	sq "github.com/Masterminds/squirrel"
)

var (
	_ LocationRepo = (*locationRepoImpl)(nil)
)

type LocationPo struct {
	ID      int64  `db:"id"`
	Address string `db:"address"`
	Country string `db:"country"`
	City    string `db:"city"`
	Area    string `db:"area"`
	ZipCode string `db:"zip_code"`
}

type LocationRepo interface {
	Find(conditions map[string]interface{}) ([]*location.Location, error)
	Save(l *location.Location) error
}

type locationRepoImpl struct {
	db *database.DB
}

func NewLocationRepo() *locationRepoImpl {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	return &locationRepoImpl{db: db}
}

func (r *locationRepoImpl) Find(conditions map[string]interface{}) ([]*location.Location, error) {
	var locations []*location.Location
	var pos []*LocationPo
	if err := r.db.Select(&pos, "SELECT * FROM location WHERE country = ?", conditions["country"]); err != nil {
		return nil, err
	}
	for _, po := range pos {
		locations = append(locations, &location.Location{
			Country: po.Country,
			City:    po.City,
			Area:    po.Area,
			ZipCode: po.ZipCode,
		})
	}
	return locations, nil
}

func (r *locationRepoImpl) Save(l *location.Location) error {
	sql, args, err := sq.Insert("locations").
		Columns("address", "country", "city", "area", "zip_code").
		Values(l.Address(), l.Country, l.City, l.Area, l.ZipCode).
		Suffix(`ON CONFLICT (address) DO UPDATE SET
			country = EXCLUDED.country,
			city = EXCLUDED.city,
			area = EXCLUDED.area,
			zip_code = EXCLUDED.zip_code
		`).
		ToSql()
	if err != nil {
		return err
	}
	if _, err := r.db.Exec(sql, args...); err != nil {
		return err
	}
	return nil
}
