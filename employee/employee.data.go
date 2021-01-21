package employee

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/maanijou/project-manager/entity"
)

// employeAPIAddress is the api address used to access employees
const employeAPIAddress = "https://employees-api.vercel.app/api/employees"

// GetEmployee is used to get an Employee by id
func GetEmployee(employeeID string) (*entity.Employee, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", employeAPIAddress, employeeID))
	if err != nil {
		log.Println("Get request Error", err)
		return nil, errors.New("Get request Error")
	}
	emp := entity.Employee{}
	err = json.NewDecoder(res.Body).Decode(&emp)
	if err != nil {
		log.Println("Error in JSON Unmarshal:", err)
		return nil, errors.New("Error in JSON Unmarshal")
	}
	return &emp, nil
}

// GetEmployees is used for getting emplooyees from external API
func GetEmployees() (*entity.Employees, error) {
	res, err := http.Get(employeAPIAddress)
	if err != nil {
		log.Println("Get request Error:", err)
		return nil, errors.New("Get request Error")
	}
	emp := entity.Employees{}
	err = json.NewDecoder(res.Body).Decode(&emp)
	if err != nil {
		log.Println("Error in JSON Unmarshal:", err)
		return nil, errors.New("Error in JSON Unmarshal")
	}
	return &emp, nil
}
