package models

import "time"

type Info struct {
	Ip       string    `bson:"ip"`
	Hostname string    `bson:"hostname"`
	Time     time.Time `bson:"time"`
	Banner   string    `bson:"banner"`
	Port     int       `bson:"port"`
	Server   string    `bson:"server"`
	WebInfo  WebInfo   `bson:"webinfo"`
}

type WebInfo struct {
	Title string   `bson:"title"`
	Tag   []string `bson:"tag"`
}
