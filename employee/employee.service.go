package employee

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maanijou/project-manager/entity"
)

// TODO handle empty employees
// GetEmployeeByID is used to get an Employee using its uuid
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eml, err := GetEmployee(vars["user_id"])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(eml)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetAllEmployees is used to get all available Employees
// TODO handle filters from query parameters
func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	page := 1
	limit := 10
	if queries["page"] != nil {
		num, err := strconv.Atoi(queries["page"][0])
		if err != nil {
			log.Println(err)
			// TODO handle redirect to same path with correct page (page=1)
		} else {
			page = num
		}
	}

	if queries["limit"] != nil {
		num, err := strconv.Atoi(queries["limit"][0])
		if err != nil {
			log.Println(err)
			// TODO handle redirect to same path with correct page (page=10)
		} else {
			limit = num
		}
	}

	if page < 1 || limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	employees, err := GetEmployees()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(employees.Data) < ((page - 1) * limit) {
		employees.Data = []*entity.Employee{}
	} else if len(employees.Data) < ((page) * limit) {
		employees.Data = employees.Data[(page-1)*limit:]
	} else {
		employees.Data = employees.Data[(page-1)*limit : page*limit]
	}
	data, err := json.Marshal(employees)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// SetupRoutes is used to setup necessary routing functions
func SetupRoutes(sm *mux.Router) {
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/employees/{user_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", GetEmployeeByID)
	getRouter.HandleFunc("/employees", GetAllEmployees)
}
