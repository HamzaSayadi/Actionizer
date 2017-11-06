package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Fullname string        `json:"fullname" bson:"fullname,omitempty"`
	ImageURL string        `json:"image_url" bson:"image_url,omitempty"`
}

type Action struct {
	ID          bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Description string        `json:"description" bson:"description,omitempty"`
}

type Task struct {
	User     User      `json:"user"`
	Action   Action    `json:"action"`
	Deadline time.Time `json:"deadline"`
}
