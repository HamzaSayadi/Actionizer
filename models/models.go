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

type Options struct {
	GetActions   bool   `short:"a" long:"getActions" description:"Show saved actions"`
	AddAction    bool   `short:"s" long:"addAction" description:"Add an actions"`
	GetUsers     bool   `short:"u" long:"getUsers" description:"Show saved users"`
	AddUser      bool   `short:"b" long:"addUser" description:"Add a user"`
	DeleteUser   bool   `short:"p" long:"deleteUser" description:"Delete a user"`
	DeleteAction bool   `short:"k" long:"deleteAction" description:"Delete an action"`
	UserName     string `short:"n" long:"username" description:"User Name to save"`
	UserImage    string `short:"i" long:"userimage" description:"User Image to save"`
	ActionDesc   string `short:"d" long:"actiondesc" description:"Action Desciption to save"`
	Config       string `short:"c" long:"config" description:"Configuration file"`
}
