package models

import "time"

type Result struct {
	Id       string    `json:"id" bson:"_id"`
	Info     string    `json:"info" bson:"info"`
	TaskId   string    `json:"task_id" bson:"task_id"`
	TaskDate time.Time `json:"task_date" bson:"task_date"`
	Ip       string    `json:"ip" bson:"ip"`
	Time     time.Time `json:"time" bson:"time"`
	Port     int       `json:"port" bson:"port"`
	VulInfo  VulInfo   `json:"vul_info" bson:"vul_info"`
}

type VulInfo struct {
	VulType  string `json:"vul_type" bson:"vul_type"`
	VulName  string `json:"vul_name" bson:"vul_name"`
	VulLevel string `json:"vul_level" bson:"vul_level"`
}
