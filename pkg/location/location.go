package location

import (
	"cake-scraper/pkg/util"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"

	_ "cake-scraper/pkg/logger"
)

var jsonPath = filepath.Join(util.ProjectRoot, "json/address.json")

var locations []*Location

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

func (l *Location) String() string {
	return l.Address()
}

// LoadLocations loads locations from json file.
func LoadLocations() []*Location {
	if locations != nil {
		return locations
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		slog.Error("failed to read json file", "error", err, "path", jsonPath)
		return nil
	}
	locations = make([]*Location, 0)
	// Save country
	const country = "Taiwan"
	location := NewLocation(country, "", "", "")
	locations = append(locations, location)
	// Iterate city
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		var city, area, zipCode string
		// Save city
		city = value.Get("city_name_en").String()
		location = NewLocation(country, city, "", "")
		locations = append(locations, location)
		// Iterate area
		value.Get("area_list").ForEach(func(key, value gjson.Result) bool {
			// Save area, zipCode
			area = value.Get("area_name_en").String()
			zipCode = value.Get("zip_code").String()
			location = NewLocation(country, city, area, zipCode)
			locations = append(locations, location)
			return true
		})
		return true
	})
	return locations
}
