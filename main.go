package main

import (
	"exam/rate"
	"log"
	"net/http"
)

func main() {
	xmlPath := "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	dbPath := "./rate.db"
	ctl := rate.NewRateController(xmlPath, dbPath)

	http.HandleFunc("/rates/latest", ctl.GetLatest)
	http.HandleFunc("/rates/", ctl.GetByTime)
	http.HandleFunc("/rates/analyze", ctl.Analyze)

	log.Printf("Starting server at port 9999\n")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}
}
