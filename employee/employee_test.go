package employee_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/maanijou/project-manager/employee"
)

func TestEmployeeStruct(t *testing.T) {
	emp := employee.Employee{}
	emp.ID = "b82522f0-8644-4c65-a552-9c6b8a9e4b6f"
	emp.FirstName = "first"
	emp.LastName = "last"
	emp.Email = "e@example.com"
	emp.Role = employee.ManagerRole
	emp.Department = employee.Engineering

	j, err := json.Marshal(emp)
	if err != nil {
		t.Errorf("Could not marshal the employee object %v", emp)
	}
	got := employee.Employee{}
	err = json.Unmarshal(j, &got)
	if err != nil {
		t.Errorf("Could not Unmarshal the employee object %v", emp)
	}
	if got != emp {
		t.Errorf("Unmarshal version of employee differs from the original %v != %v", got, emp)
	}

	emp.Role = "other"
	j, err = json.Marshal(emp)
	err = json.Unmarshal(j, &got)
	if err == nil {
		t.Errorf("Expected error checking role: %v", err)
	}
	emp.Role = employee.ManagerRole
	emp.Department = "other"
	j, err = json.Marshal(emp)
	err = json.Unmarshal(j, &got)

	if err == nil {
		t.Errorf("Expected error checking Department: %v", err)
	}

}

func TestGetEmployeeData(t *testing.T) {
	var expect *employee.Employee = &employee.Employee{}
	expect.ID = "b82522f0-8644-4c65-a552-9c6b8a9e4b6f"
	expect.FirstName = ""
	expect.LastName = "Simenot"
	expect.Email = "simenot@acme.com"
	expect.Role = employee.ManagerRole
	expect.Department = employee.Marketing

	got, err := employee.GetEmployee("b82522f0-8644-4c65-a552-9c6b8a9e4b6f")
	if err != nil {
		t.Errorf("Error getting emp from External server")
	}
	if *got != *expect {
		t.Errorf("Error getting expected employee %v != %v", got, expect)
	}
}

func TestGetEmployees(t *testing.T) {
	emps, err := employee.GetEmployees()
	if err != nil {
		t.Errorf("Error in getting employees %v", err)
	}
	if len(emps.Data) != 100 {
		t.Errorf("Error in getting employees, expected 100 employees got %d", len(emps.Data))
	}

}

func TestGetEmployeeByIDHandler(t *testing.T) {
	sm := mux.NewRouter().StrictSlash(true) // ignoring trailing slash
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/employees/{user_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", employee.GetEmployeeByID)
	req, err := http.NewRequest("GET", "/employees/a0d5e87a-af04-473d-b1f5-3105bbf986c8", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.Handler(getRouter)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":"a0d5e87a-af04-473d-b1f5-3105bbf986c8","first_name":"Celia","last_name":"Ladbrook","email":"celia.ladbrook@acme.com","department":"sales","role":"employee"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetEmployeesHandler(t *testing.T) {
	sm := mux.NewRouter().StrictSlash(true) // ignoring trailing slash
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/employees", employee.GetAllEmployees)

	// table tests
	scenarios := []struct {
		page         string
		limit        string
		expect       int
		expectAnswer string
	}{
		{page: "4", limit: "2", expect: 2},
		{page: "12", limit: "9", expect: 1, expectAnswer: "Jackie"},
		{page: "12", limit: "9", expect: 1},
		{page: "", limit: "2", expect: 2},
		{page: "", limit: "", expect: 10},
	}
	for i, s := range scenarios {
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("/employees?page=%s&limit=%s", s.page, s.limit),
			nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.Handler(getRouter)
		handler.ServeHTTP(rr, req)
		_, errPage := strconv.Atoi(s.page)
		_, errLimit := strconv.Atoi(s.limit)
		if status := rr.Code; errPage == nil &&
			errLimit == nil &&
			status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
			continue
		}

		// Check the response body is what we expect.
		var employees employee.Employees
		err = json.Unmarshal(rr.Body.Bytes(), &employees)
		if err != nil {
			t.Errorf("Error Unmarshaling json response in Scenario %d:", i)
		}

		if len(employees.Data) != s.expect {
			t.Errorf("handler returned unexpected body: got %v employees want %v employees",
				len(employees.Data), s.expect)
		}
		if len(s.expectAnswer) > 0 && len(employees.Data) > 0 && !strings.EqualFold(employees.Data[0].FirstName, s.expectAnswer) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				employees.Data[0].FirstName, s.expectAnswer)
		}
	}

}
