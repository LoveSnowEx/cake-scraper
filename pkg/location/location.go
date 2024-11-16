package location

import (
	"cake-scraper/pkg/util"
	"strings"
)

type Location struct {
	Country string
	City    string
	Area    string
	ZipCode string
}

func NewLocation(country, city, area, zipCode string) *Location {
	return &Location{
		Country: country,
		City:    city,
		Area:    area,
		ZipCode: zipCode,
	}
}

func (l *Location) Address() string {
	return strings.Join(
		util.Filter(
			[]string{l.Area, l.City, l.Country},
			func(s string) bool { return s != "" },
		),
		", ",
	)
}
