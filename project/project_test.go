package project_test

import (
	"encoding/json"
	"testing"

	"github.com/maanijou/project-manager/employee"
	"github.com/maanijou/project-manager/project"
)

func TestProjectStruct(t *testing.T) {
	emp := employee.Employee{}
	emp.ID = "b82522f0-8644-4c65-a552-9c6b8a9e4b6f"
	emp.FirstName = "first"
	emp.LastName = "last"
	emp.Email = "e@example.com"
	emp.Role = employee.ManagerRole
	emp.Department = employee.Engineering

	prj := project.Project{}
	prj.ID = 1
	prj.Name = "project"
	prj.Owner = emp
	prj.Participants = []employee.Employee{emp}
	prj.Progress = 0
	prj.State = project.Failed

	// TODO check progress and state together
	j, err := json.Marshal(prj)
	if err != nil {
		t.Errorf("Could not marshal the project object %v", prj)
	}
	got := project.Project{}
	err = json.Unmarshal(j, &got)
	if err != nil {
		t.Errorf("Could not Unmarshal the project object %v", prj)
	}

}
