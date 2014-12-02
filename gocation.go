package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var displayURL bool = false
var db *sql.DB

func init() {
	flag.BoolVar(&displayURL, "db-url", false, "output database url and exit")

	var err error
	db, err = sql.Open(DB_Driver, config.DB_URL().String())
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {
	flag.Parse()

	if displayURL {
		fmt.Printf(config.DB_URL().String())
		os.Exit(0)
	}
}

func main() {
	parseFlags()

	ts, err := time.Parse(time.RFC3339, "2014-12-01T21:16:32-08:00")
	if err != nil {
		log.Fatal(err)
	}

	l := &Location{
		Longitude: -122.6764,
		Latitude:  45.5165,
		Accuracy:  10,
		Timestamp: &ts}
	l.Save()
	fmt.Printf("saved id: %d\n", l.Id)

	l.UpdateGeo()

	x := LoadLocation(int64(1))
	fmt.Println(x.ToJSON())

	/*
		fmt.Printf("point - %s\n", x.Point())
		fmt.Printf("polygon - %s\n", x.Polygon())
	*/

	defer db.Close()
}
