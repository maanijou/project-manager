package employee

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// employeAPIAddress is the api address used to access employees
const employeAPIAddress = "https://employees-api.vercel.app/api/employees"

// GetEmployee is used to get an Employee by id
func GetEmployee(employeeID string) (*Employee, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", employeAPIAddress, employeeID))
	if err != nil {
		log.Println("Get request Error", err)
		return nil, errors.New("Get request Error")
	}
	emp := Employee{}
	err = json.NewDecoder(res.Body).Decode(&emp)
	if err != nil {
		log.Println("Error in JSON Unmarshal:", err)
		return nil, errors.New("Error in JSON Unmarshal")
	}
	return &emp, nil
}

// GetEmployees is used for getting emplooyees from external API
func GetEmployees() (*Employees, error) {
	res, err := http.Get(employeAPIAddress)
	if err != nil {
		log.Println("Get request Error:", err)
		return nil, errors.New("Get request Error")
	}
	emp := Employees{}
	err = json.NewDecoder(res.Body).Decode(&emp)
	if err != nil {
		log.Println("Error in JSON Unmarshal:", err)
		return nil, errors.New("Error in JSON Unmarshal")
	}
	return &emp, nil
}
