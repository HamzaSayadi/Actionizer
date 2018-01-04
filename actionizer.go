package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/knadh/jsonconfig"
	"github.com/syd7/actionizer/cli"
	"github.com/syd7/actionizer/database"
	"github.com/syd7/actionizer/models"
	"github.com/syd7/actionizer/server"
	"github.com/syd7/actionizer/utils"
)

type configuration struct {
	Database       database.Config `json:"database"`
	ActionDuration utils.Duration  `json:"action_duration"`
	HTTPListen     string          `json:"http_listen"`
	HTTPPort       int             `json:"http_port"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	//parse cli options
	var opts models.Options
	_, erro := flags.Parse(&opts)
	fmt.Println(opts)

	if erro != nil {
		panic(erro)
	}

	configFile := opts.Config
	if configFile == "" {
		configFile = "actionizer.json"
	}
	// parse and load json config
	var config configuration
	err := jsonconfig.Load(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}

	db, err := database.Connect(config.Database)
	if err != nil {
		log.Fatalf("Cannot connect to database server: %v\n", err)
	}

	fmt.Println(opts.GetUsers)
	if opts.GetActions {
		cli.ShowActions(db)
		return
	}
	if opts.GetUsers {
		cli.ShowUsers(db)
		return
	}
	if opts.AddUser {
		cli.AddUser(db, opts.UserName, opts.UserImage)
		return
	}
	if opts.AddAction {
		cli.AddAction(db, opts.ActionDesc)
		return
	}
	if opts.DeleteUser {
		cli.DeleteUser(db, opts.UserName)
	}
	if opts.DeleteAction {
		cli.DeleteAction(db, opts.ActionDesc)
	}

	//Choose an action if there none
	_, err = db.CurrentTask()
	if err != nil {
		log.Printf("No task found, creating new one\n")
		db.NewRandomTask(config.ActionDuration)
	}

	go func() {
		c := time.Tick(time.Duration(config.ActionDuration))
		for _ = range c {
			db.NewRandomTask(config.ActionDuration)
			log.Printf("New task created")
		}
	}()

	server := server.Server{
		Host: config.HTTPListen,
		Port: config.HTTPPort,
		DB:   db,
	}

	log.Printf("Listening on %s:%d\n", config.HTTPListen, config.HTTPPort)
	err = server.Run()
	log.Fatal(err)
}
