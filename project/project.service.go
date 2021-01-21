package project

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// GetProjectByID is used as a route to get a project by its id
func GetProjectByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	project, err := getProject(projectID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			w.Header().Add("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf(`{"error": "Project not found"}`)))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	j, err := json.Marshal(project)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetProjectList is used to get list of all projects
// Can accept page and limit query parameters (positives)
func GetProjectList(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	page := 1
	limit := 10
	if queries["page"] != nil {
		num, err := strconv.Atoi(queries["page"][0])
		if err != nil {
			page = 1
		}
		page = num
	}

	if queries["limit"] != nil {
		num, err := strconv.Atoi(queries["limit"][0])
		if err != nil {
			limit = 10
		}
		limit = num
	}

	projects, err := getProjects(page, limit)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(projects)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func deleteProjectRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = removeProject(projectID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

// SetupRoutes is used to setup routing for project services
func SetupRoutes(sm *mux.Router) {
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	getRouter.HandleFunc("/projects/", GetProjectList)
	getRouter.HandleFunc("/projects/{id:[0-9]+}", GetProjectByID)
	postRouter.HandleFunc("/projects", insertProjectRoute)
	deleteRouter.HandleFunc("/projects/{id:[0-9]+}", deleteProjectRoute)
}
