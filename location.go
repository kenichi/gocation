package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

const (
	locationInsertSQL = `
		INSERT INTO locations(
			"longitude",
			"latitude",
			"accuracy",
			"speed",
			"bearing",
			"timestamp",
			"altitude",
			"vertical_accuracy",
			"battery",
			"charging")
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`

	locationSelectSQL = `
		SELECT
			"longitude",
			"latitude",
			"accuracy",
			"speed",
			"bearing",
			"timestamp",
			"altitude",
			"vertical_accuracy",
			"battery",
			"charging"
		FROM locations
		WHERE "id"=$1
		LIMIT 1`

	locationUpdateGeometrySQL = `
		UPDATE locations SET
			"point"=ST_GeogFromWKB(ST_MakePoint("longitude","latitude")),
			"polygon"=ST_GeomFromEWKB(ST_Buffer(ST_GeogFromWKB(ST_MakePoint("longitude","latitude")), "accuracy"))
		WHERE "id"=$1`

	/*
		locationSelectPointSQL = `
			SELECT ST_AsGeoJSON("point")
			FROM locations
			WHERE "id"=$1`

		locationSelectPolygonSQL = `
			SELECT ST_AsGeoJSON("polygon")
			FROM locations
			WHERE "id"=$1`
	*/
)

type Location struct {
	Id               int64      `json:"id"`
	Longitude        float64    `json:"longitude"`
	Latitude         float64    `json:"latitude"`
	Accuracy         float64    `json:"accuracy"`
	Speed            int64      `json:"speed"`
	Bearing          int64      `json:"bearing"`
	Timestamp        *time.Time `json:"timestamp"`
	Altitude         int64      `json:"altitude"`
	VerticalAccuracy int64      `json:"verticalAccuracy"`
	Battery          int64      `json:"battery"`
	Charging         bool       `json:"charging"`
}

var locationInsertStmt *sql.Stmt
var locationSelectStmt *sql.Stmt
var locationUpdateGeometryStmt *sql.Stmt

/*
var locationSelectPointStmt *sql.Stmt
var locationSelectPolygonStmt *sql.Stmt
*/

func init() {
	var err error

	locationInsertStmt, err = db.Prepare(locationInsertSQL)
	if err != nil {
		log.Fatal(err)
	}

	locationSelectStmt, err = db.Prepare(locationSelectSQL)
	if err != nil {
		log.Fatal(err)
	}

	locationUpdateGeometryStmt, err = db.Prepare(locationUpdateGeometrySQL)
	if err != nil {
		log.Fatal(err)
	}

	/*
		locationSelectPointStmt, err = db.Prepare(locationSelectPointSQL)
		if err != nil {
			log.Fatal(err)
		}

		locationSelectPolygonStmt, err = db.Prepare(locationSelectPolygonSQL)
		if err != nil {
			log.Fatal(err)
		}
	*/
}

func LoadLocation(id int64) *Location {
	l := &Location{}
	if err := locationSelectStmt.QueryRow(id).Scan(
		&l.Longitude,
		&l.Latitude,
		&l.Accuracy,
		&l.Speed,
		&l.Bearing,
		&l.Timestamp,
		&l.Altitude,
		&l.VerticalAccuracy,
		&l.Battery,
		&l.Charging); err != nil {
		log.Fatal(err)
	}
	l.Id = id
	lts := l.Timestamp.Local()
	l.Timestamp = &lts
	return l
}

func (l *Location) Save() int64 {
	locationInsertStmt.QueryRow(
		l.Longitude,
		l.Latitude,
		l.Accuracy,
		l.Speed,
		l.Bearing,
		l.Timestamp.UTC(),
		l.Altitude,
		l.VerticalAccuracy,
		l.Battery,
		l.Charging).Scan(&l.Id)
	return l.Id
}

func (l *Location) UpdateGeo() {
	locationUpdateGeometryStmt.Exec(l.Id)
}

/*
func (l *Location) Point() string {
	var point string
	if err := locationSelectPointStmt.QueryRow(l.Id).Scan(&point); err != nil {
		log.Fatal(err)
	}
	return point
}

func (l *Location) Polygon() string {
	var polygon string
	if err := locationSelectPolygonStmt.QueryRow(l.Id).Scan(&polygon); err != nil {
		log.Fatal(err)
	}
	return polygon
}
*/

func (l *Location) ToJSON() string {
	jsonBytes, err := json.Marshal(l)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonBytes[:])
}
