package models

import "time"

type Plugin struct {
	Id       string    `json:"id" bson:"_id"`
	Count    int       `bson:"count" json:"count"`
	AddTime  time.Time `bson:"add_time" json:"add_time"`
	Info     string    `bson:"info" json:"info"`
	Name     string    `bson:"name" json:"name"`
	Keyword  string    `bson:"keyword" json:"keyword"`
	Level    string    `bson:"level" json:"level"`
	Url      string    `bson:"url" json:"url"`
	Author   string    `bson:"author" json:"author"`
	Filename string    `bson:"filename" json:"filename"`
	Source   int       `bson:"source" json:"source"`
	Type     string    `bson:"type" json:"type"`
}

type VulType struct {
	Type  string `bson:"type" json:"type"`
	Count int    `bson:"value" json:"value"`
}
