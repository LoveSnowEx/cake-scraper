package main

import (
	"cake-scraper/pkg/repo/locationrepo"
	"fmt"
)

func main() {
	if err := locationrepo.NewLocationRepo().Init(); err != nil {
		fmt.Println(err)
		return
	}
}
