package models

import "time"

type Task struct {
	Id        string          `json:"id" bson:"_id"`
	Title     string          `json:"title" bson:"title"`
	Status    int             `json:"status" bson:"status"`
	IsUpdate  bool            `json:"isupdate" bson:"isupdate"`
	Time      time.Time       `json:"time" bson:"time"`
	Query     string          `json:"query" bson:"query"`
	Plan      int             `json:"plan" bson:"plan"`
	Plugin    string          `json:"plugin" bson:"plugin"`
	Condition string          `json:"condition" bson:"condition"`
	Target    [][]interface{} `json:"target" bson:"target"`
}

type InTask struct {
	Title     string          `json:"title" bson:"title"`
	Status    int             `json:"status" bson:"status"`
	IsUpdate  bool            `json:"isupdate" bson:"isupdate"`
	Time      time.Time       `json:"time" bson:"time"`
	Query     string          `json:"query" bson:"query"`
	Plan      int             `json:"plan" bson:"plan"`
	Plugin    string          `json:"plugin" bson:"plugin"`
	Condition string          `json:"condition" bson:"condition"`
	Target    [][]interface{} `json:"target" bson:"target"`
}

type AllData struct {
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Hostname string `json:"hostname"`
	VulLevel string `json:"vul_Level"`
	Info     string `json:"info"`
	VulName  string `json:"vul_Name"`
	Title    string `json:"title"`
	Time     string `json:"time"`
	LastScan string `json:"lastscan"`
}
