package location

import (
	"cake-scraper/pkg/util"
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
	if l == nil {
		return ""
	}
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
	util.PanicError(err)
	locations = make([]*Location, 0)
	// Iterate country
	gjson.ParseBytes(data).ForEach(func(key, country_node gjson.Result) bool {
		country := country_node.Get("country_name_en").String()
		location := NewLocation(country, "", "", "")
		locations = append(locations, location)
		// Iterate city
		country_node.Get("city_list").ForEach(func(key, city_node gjson.Result) bool {
			city := city_node.Get("city_name_en").String()
			location = NewLocation(country, city, "", "")
			locations = append(locations, location)
			// Iterate area
			city_node.Get("area_list").ForEach(func(key, area_node gjson.Result) bool {
				area := area_node.Get("area_name_en").String()
				zipCode := area_node.Get("zip_code").String()
				location = NewLocation(country, city, area, zipCode)
				locations = append(locations, location)
				return true
			})
			return true
		})
		return true
	})
	return locations
}
