package models

type Statistics struct {
	Id   string         `json:"id" bson:"_id"`
	Date string         `json:"date" bson:"date"`
	Info StatisticsInfo `json:"info" bson:"info"`
}

type StatisticsInfo struct {
	Add    int `json:"add" bson:"add"`
	Update int `json:"update" bson:"update"`
	Delete int `json:"delete" bson:"delete"`
}
