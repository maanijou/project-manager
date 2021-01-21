package project

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func updateProjectRoute(w http.ResponseWriter, r *http.Request) {
	data, eerr := ioutil.ReadAll(r.Body)
	log.Println(eerr)

	defer r.Body.Close()
	log.Println("data", string(data))
	var project entity.Project
	// err := json.NewDecoder(r.Body).Decode(&project)
	err := json.Unmarshal(data, &project)
	if err != nil {
		log.Println("Error in json decoder:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("This is the project %v", project)
	projectID, err := updateProject(project)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf(`{"project_id": %d}`, projectID)))
}

func addParticipantsRoute(w http.ResponseWriter, r *http.Request) {
	// return
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["project_id"])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var participants []entity.Employee
	err = json.NewDecoder(r.Body).Decode(&participants)
	count, err := addParticipant(projectID, participants)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			w.Header().Add("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf(`{"error": "Project not found"}`)))
		} else if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {

		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if count > 0 {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"result": "%d participants added out of %d"}`, count, len(participants))))
	} else {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"result": "%d participants added out of %d"}`, count, len(participants))))
	}
	return
}

func removeParticipantsRoute(w http.ResponseWriter, r *http.Request) {
	// return
	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["project_id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	participantID := vars["participant_id"]
	if len(participantID) < 1 {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = removeParticipant(projectID, participantID)
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"result": "Participant not exist or cannot be removed."}`)))
	}
	w.WriteHeader(http.StatusNoContent)
}

// SetupRoutes is used to setup routing for project services
func SetupRoutes(sm *mux.Router) {
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	getRouter.HandleFunc("/projects/", GetProjectList)
	getRouter.HandleFunc("/projects/{id:[0-9]+}", GetProjectByID)
	postRouter.HandleFunc("/projects", insertProjectRoute)
	postRouter.HandleFunc("/projects/{project_id:[0-9]+}/participants", addParticipantsRoute)
	deleteRouter.HandleFunc("/projects/{id:[0-9]+}", deleteProjectRoute)
	deleteRouter.HandleFunc(
		"/projects/{project_id:[0-9]+}/participants/{participant_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}",
		removeParticipantsRoute)
	putRouter.HandleFunc("/projects", updateProjectRoute)
}
