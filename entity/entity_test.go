package entity_test

import (
	"encoding/json"
	"testing"

	"github.com/maanijou/project-manager/entity"
)

func TestEmployeeStruct(t *testing.T) {
	emp := entity.Employee{}
	emp.ID = "b82522f0-8644-4c65-a552-9c6b8a9e4b6f"
	emp.FirstName = "first"
	emp.LastName = "last"
	emp.Email = "e@example.com"
	emp.Role = entity.ManagerRole
	emp.Department = entity.Engineering

	j, err := json.Marshal(emp)
	if err != nil {
		t.Errorf("Could not marshal the employee object %v", emp)
	}
	got := entity.Employee{}
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
	emp.Role = entity.ManagerRole
	emp.Department = "other"
	j, err = json.Marshal(emp)
	err = json.Unmarshal(j, &got)

	if err == nil {
		t.Errorf("Expected error checking Department: %v", err)
	}

}
