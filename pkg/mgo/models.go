package mgo

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type ListRes []ListData

type ListData struct {
	Items []bson.M `json:"items"`
	Total []Total  `json:"total"`
}

type Total struct {
	Total int `json:"total"`
}

func (l ListRes) Convert(data interface{}) error {
	if len(l) == 0 {
		return nil
	}
	bytes, err := json.Marshal(l[0].Items)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, data)
}

func (l ListRes) GetTotal() int {
	if len(l) == 0 {
		return 0
	}
	return l[0].GetTotal()
}

func (l *ListData) GetTotal() int {
	if l == nil || len(l.Total) == 0 {
		return 0
	}
	return l.Total[0].Total
}
