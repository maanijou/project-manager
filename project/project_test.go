package project_test

import (
	"encoding/json"
	"testing"

	"github.com/maanijou/project-manager/entity"
)

func TestProjectStruct(t *testing.T) {
	emp := entity.Employee{}
	emp.ID = "b82522f0-8644-4c65-a552-9c6b8a9e4b6f"
	emp.FirstName = "first"
	emp.LastName = "last"
	emp.Email = "e@example.com"
	emp.Role = entity.ManagerRole
	emp.Department = entity.Engineering

	prj := entity.Project{}
	prj.ID = 1
	prj.Name = "project"
	prj.Owner = emp
	prj.Participants = []entity.Employee{emp}
	prj.Progress = 0
	prj.State = entity.Failed

	// TODO check progress and state together
	j, err := json.Marshal(prj)
	if err != nil {
		t.Errorf("Could not marshal the project object %v", prj)
	}
	got := entity.Project{}
	err = json.Unmarshal(j, &got)
	if err != nil {
		t.Errorf("Could not Unmarshal the project object %v", prj)
	}

}
