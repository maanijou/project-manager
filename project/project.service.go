package project

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/maanijou/project-manager/entity"
)

// insertProjectRoute is a rout to handle project insertion
func insertProjectRoute(w http.ResponseWriter, r *http.Request) {
	var project entity.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	projectID, err := insertProject(project)
	if err != nil {
		log.Println("Error in inserting project:", err)
		w.WriteHeader(http.StatusBadRequest)
		if strings.Contains(err.Error(), "duplicate key value") {
			w.Write([]byte(fmt.Sprintf(`{"error": %s}`, "Project name already exists")))
		}
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"project_id": %d}`, projectID)))
}

// SetupRoutes is used to setup routing for project services
func SetupRoutes(sm *mux.Router) {
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/projects", insertProjectRoute)
}
