package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/syd7/actionizer/database"
)

func getCurrentTask(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		task, err := db.CurrentTask()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(task)
	}
}

func getUsers(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		users := db.AllUsers()
		json.NewEncoder(w).Encode(users)
	}
}

func getActions(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		actions := db.AllActions()
		json.NewEncoder(w).Encode(actions)
	}
}

func getTasks(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tasks := db.AllTasks()
		json.NewEncoder(w).Encode(tasks)
	}
}

type Server struct {
	Host string
	Port int
	DB   *database.Database
}

func (s Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	http.HandleFunc("/current_task", getCurrentTask(s.DB))
	http.HandleFunc("/users", getUsers(s.DB))
	http.HandleFunc("/actions", getActions(s.DB))
	http.HandleFunc("/tasks", getTasks(s.DB))

	http.Handle("/", http.FileServer(http.Dir("public/")))

	return http.ListenAndServe(addr, nil)
}
