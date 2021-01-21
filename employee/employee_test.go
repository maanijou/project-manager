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
}
