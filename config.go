package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

const (
	File        = "config.json"
	DB_Driver   = "postgres"
	DB_RawQuery = "sslmode=disable"
)

type Config struct {
	DB struct {
		Name string `json:"name"`
		Host string `json:"host"`
		Port int    `json:"port"`
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"db"`
}

var config Config

func init() {
	if body, err := ioutil.ReadFile(File); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(body, &config); err != nil {
			log.Println("error loading config.json")
			panic(err)
		}
	}
}

func (c *Config) DB_URL() *url.URL {
	return &url.URL{
		Scheme:   DB_Driver,
		User:     url.UserPassword(c.DB.User, c.DB.Pass),
		Host:     fmt.Sprintf("%s:%d", c.DB.Host, c.DB.Port),
		Path:     c.DB.Name,
		RawQuery: DB_RawQuery}
}
