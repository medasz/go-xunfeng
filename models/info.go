package models

import "time"

type Info struct {
	Ip       string    `bson:"ip"`
	Hostname string    `bson:"hostname,omitempty"`
	Time     time.Time `bson:"time"`
	Banner   string    `bson:"banner"`
	Port     int       `bson:"port"`
	Server   string    `bson:"server,omitempty"`
	WebInfo  WebInfo   `bson:"webinfo,omitempty"`
}

type WebInfo struct {
	Title string   `bson:"title,omitempty"`
	Tag   []string `bson:"tag,omitempty"`
}

type TypeCount struct {
	ServerName string `bson:"server_name" json:"server_name"`
	Count      int    `bson:"count" json:"count"`
}
