package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/think-it-labs/actionizer/database"
	"github.com/think-it-labs/actionizer/server"
	"github.com/think-it-labs/actionizer/utils"

	"github.com/knadh/jsonconfig"
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
	configFile := flag.String("config", "config/config.json", "Configuration file")

	// parse and load json config
	var config configuration
	err := jsonconfig.Load(*configFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}

	db, err := database.Connect(config.Database)
	if err != nil {
		log.Fatalf("Cannot connect to database server: %v\n", err)
	}

	// Choose an action if there none
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
