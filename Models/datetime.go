package models

type Datetime struct {
	Day  string `json:"dia" bson:"day,omitempty"`
	Time string `json:"hora" bson:"time,omitempty"`
}
