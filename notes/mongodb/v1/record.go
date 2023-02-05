package main

import "time"

type Record struct {
	ID          string    `bson:"_id,omitempty"`
	Name        string    `bson:"name"`
	CreatedTime time.Time `bson:"createdTime"`
	Count       int64     `bson:"count"`
	IsTest      bool      `bson:"isTest"`
}
