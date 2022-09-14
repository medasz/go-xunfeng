package models

import "time"

type Heartbeat struct {
	Name   string    `json:"name" bson:"name"`
	UpTime time.Time `json:"up_time" bson:"up_time"`
	Value  float64   `json:"value" bson:"value"`
}
