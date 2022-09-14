package models

import "time"

type History struct {
	Info
	DelTime time.Time `bson:"del_time"`
	Type    string    `bson:"type"`
}
