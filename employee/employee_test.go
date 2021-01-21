package employee_test

import (
	"encoding/json"
	"testing"

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

func TestGetEmployeeByIDHandler(t *testing.T) {
	emps, err := employee.GetEmployees()
	if err != nil {
		t.Errorf("Error in getting employees %v", err)
	}
	if len(emps.Data) != 100 {
		t.Errorf("Error in getting employees, expected 100 employees got %d", len(emps.Data))
	}

}
