package database

import (
	"math/rand"
	"time"

	"github.com/syd7/actionizer/models"
	"github.com/think-it-labs/actionizer/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Config struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Database struct {
	mongodb *mgo.Database
}

type taskAssociation struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID   bson.ObjectId `json:"user_id" bson:"user_id,omitempty"`
	ActionID bson.ObjectId `json:"action_id" bson:"action_id,omitempty"`
	Deadline time.Time     `json:"deadline" bson:"deadline,omitempty"`
}

func Connect(config Config) (*Database, error) {
	session, err := mgo.Dial(config.Host)
	if err != nil {
		return nil, err
	}

	db := session.DB(config.Name)

	// If User is not empty we do an auth
	if config.User != "" {
		err := db.Login(config.User, config.Password)
		if err != nil {
			return nil, err
		}
	}
	return &Database{db}, nil
}

func (db *Database) ucol() *mgo.Collection {
	return db.mongodb.C("users")
}

func (db *Database) acol() *mgo.Collection {
	return db.mongodb.C("actions")
}

func (db *Database) tcol() *mgo.Collection {
	return db.mongodb.C("tasks")
}

func (db *Database) CurrentTask() (*models.Task, error) {
	var taskAssociation taskAssociation
	now := time.Now().UTC()
	err := db.tcol().Find(
		bson.M{
			"deadline": bson.M{
				"$gt": now,
			},
		}).One(&taskAssociation)
	if err != nil {
		return nil, err
	}

	return db.GetTask(taskAssociation)
}

func (db *Database) GetTask(t taskAssociation) (*models.Task, error) {
	var task models.Task

	// Query for the action
	err := db.acol().Find(
		bson.M{
			"_id": t.ActionID,
		}).One(&task.Action)
	if err != nil {
		return nil, err
	}

	// Query for the user
	err = db.ucol().Find(
		bson.M{
			"_id": t.UserID,
		}).One(&task.User)
	if err != nil {
		return nil, err
	}

	// Set the deadline
	task.Deadline = t.Deadline

	return &task, nil
}

func (db *Database) AllUsers() []models.User {
	var users []models.User
	db.ucol().Find(nil).All(&users)
	return users
}

func (db *Database) AllActions() []models.Action {
	var actions []models.Action
	db.acol().Find(nil).All(&actions)
	return actions
}

func (db *Database) AllTasks() []models.Task {
	var tasks []models.Task
	var tasksAssoc []taskAssociation
	db.tcol().Find(nil).All(&tasksAssoc)
	for _, taskAssoc := range tasksAssoc {
		task, err := db.GetTask(taskAssoc)
		if err == nil {
			tasks = append(tasks, *task)
		}
	}
	return tasks
}

func (db *Database) NewRandomTask(taskDuration utils.Duration) (models.Task, error) {
	users := db.AllUsers()
	actions := db.AllActions()
	choosenUser := users[rand.Intn(len(users))]
	choosenAction := actions[rand.Intn(len(actions))]
	deadline := time.Now().UTC().Add(time.Duration(taskDuration))

	return db.AffectActionToUser(choosenUser, choosenAction, deadline)
}

func (db *Database) InsertUser(u models.User) error {
	ucol := db.ucol()
	err := ucol.Insert(&u)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteUser(name string) error {

	err := db.ucol().Remove(
		bson.M{
			"fullname": name,
		})
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) InsertAction(a models.Action) error {
	acol := db.acol()
	err := acol.Insert(&a)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteAction(description string) error {

	err := db.acol().Remove(
		bson.M{
			"description": description,
		})
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) AffectActionToUser(u models.User, a models.Action, d time.Time) (models.Task, error) {
	taskAssoc := taskAssociation{
		UserID:   u.ID,
		ActionID: a.ID,
		Deadline: d,
	}

	err := db.tcol().Insert(&taskAssoc)
	if err != nil {
		return models.Task{}, err
	}
	return models.Task{User: u, Action: a, Deadline: d}, nil
}
