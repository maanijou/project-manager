package employee

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO handle empty employees
// GetEmployeeByID is used to get an Employee using its uuid
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eml, err := GetEmployee(vars["user_id"])

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

// SetupRoutes is used to setup necessary routing functions
func SetupRoutes(sm *mux.Router) {
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/employees/{user_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", GetEmployeeByID)
}
