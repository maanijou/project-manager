package entity

import (
	"encoding/json"
	"errors"
	"log"
)

// DepartmentType type (int)
type DepartmentType string

const (
	// Sales is "sales"
	Sales DepartmentType = "sales"
	// Engineering is "engineering"
	Engineering DepartmentType = "engineering"
	// Marketing is "marketing"
	Marketing DepartmentType = "marketing"
)

// RoleType type string
type RoleType string

const (
	// ManagerRole is "manager"
	ManagerRole RoleType = "manager"
	// EmployeeRole is "employee"
	EmployeeRole RoleType = "employee"
)

var roleName = map[RoleType]string{
	ManagerRole:  "manager",
	EmployeeRole: "employee",
}

func (x RoleType) String() string {
	return roleName[x]
}

// Employee is a struct
type Employee struct {
	ID         string         `json:"id"`
	FirstName  string         `json:"first_name,omitempty" `
	LastName   string         `json:"last_name,omitempty"`
	Email      string         `json:"email,omitempty" `
	Department DepartmentType `json:"department,omitempty" `
	Role       RoleType       `json:"role,omitempty"`
}

// Employees is an array of type Employee
type Employees struct {
	Data []*Employee `json:"data"`
}

// UnmarshalJSON is
func (s *Employee) UnmarshalJSON(data []byte) error {
	// Define a secondary type so that we don't end up with a recursive call to json.Unmarshal
	type Aux Employee
	var a *Aux = (*Aux)(s)
	err := json.Unmarshal(data, &a)
	if err != nil {
		log.Println("Error in UnmarshalJSON", err)
		return err
	}

	// Validate the valid enum values
	switch s.Department {
	case Engineering, Marketing, Sales, "":
	default:
		log.Println("Error in UnmarshalJSON, invalid value for Department")
		return errors.New("invalid value for Department")
	}

	// Validate the valid enum values
	switch s.Role {
	case ManagerRole, EmployeeRole, "":
		return nil
	default:
		log.Println("Error UnmarshalJSON, invalid value for Role")
		return errors.New("invalid value for Role")
	}
}
