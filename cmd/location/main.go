package main

import (
	"cake-scraper/pkg/location"
	"cake-scraper/pkg/repo/locationrepo"
	"log"
	"os"

	"github.com/tidwall/gjson"
)

const jsonPath = "json/address.json"

func main() {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalln(err)
	}
	repo := locationrepo.NewLocationRepo()

	// Save country
	const country = "Taiwan"
	if err := repo.Save(location.NewLocation(country, "", "", "")); err != nil {
		log.Fatalln(err)
	}
	// Iterate city
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		var city, area, zipCode string
		// Save city
		city = value.Get("city_name_en").String()
		if err := repo.Save(location.NewLocation(country, city, area, zipCode)); err != nil {
			log.Fatalln(err)
		}
		// Iterate area
		value.Get("area_list").ForEach(func(key, value gjson.Result) bool {
			// Save area, zipCode
			area = value.Get("area_name_en").String()
			zipCode = value.Get("zip_code").String()
			if err := repo.Save(location.NewLocation(country, city, area, zipCode)); err != nil {
				log.Fatalln(err)
			}
			return true
		})
		return true
	})
}
