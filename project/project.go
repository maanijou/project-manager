package project

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/maanijou/project-manager/employee"
)

// StateType is string
type StateType string

const (
	// Planned is "planned" 0
	Planned StateType = "planned"
	// Active is "active" 1
	Active StateType = "active"
	// Done is "done" 2
	Done StateType = "done"
	// Failed is "failed" 3
	Failed StateType = "failed"
)

// StateName is used for getting string value of int states
var StateName = map[int]StateType{
	0: "planned",
	1: "active",
	2: "done",
	3: "failed",
}

var stateNum = map[StateType]int{
	Planned: 0,
	Active:  1,
	Done:    2,
	Failed:  3,
}

// Int is used to get int value of states
func (x StateType) Int() int {
	return stateNum[x]
}

// Project is a struct for handling projects
type Project struct {
	ID           uint                `json:"id"`
	Name         string              `json:"name"`
	Owner        employee.Employee   `json:"owner"`
	State        StateType           `json:"state"`
	Progress     int                 `json:"progress,omitempty"`
	Participants []employee.Employee `json:"participants" `
}

// UnmarshalJSON for Project
func (s *Project) UnmarshalJSON(data []byte) error {
	// Define a secondary type so that we don't end up with a recursive call to json.Unmarshal
	type Aux Project
	var a *Aux = (*Aux)(s)
	err := json.Unmarshal(data, &a)
	if err != nil {
		log.Print("Error in here", err)
		return err
	}

	// TODO add check for progess to be set only for active projects

	// Validate the valid enum values
	switch s.State {
	case Active, Planned, Failed, Done:
	default:
		log.Println("invalid value for state")
		return errors.New("invalid value for state")
	}
	return nil
}

// TODO add compare method for project in here
