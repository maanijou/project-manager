package employee

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maanijou/project-manager/database"
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

// InsertEmployee for inserting Employees
func InsertEmployee(employee *entity.Employee) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	employeeID, err := uuid.FromString(employee.ID)
	if err != nil {
		return 0, err
	}
	stm, err := database.DbConn.Prepare(`INSERT INTO employee(
		id)
		VALUES ($1);
	`)
	if err != nil {
		log.Printf("Error inserting employee %v\n", err)
		return 0, nil
	}
	result, err := stm.ExecContext(ctx, employeeID)
	if err != nil {
		log.Println("pg error:", err.Error())
		return 0, err
	}
	num, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(num), nil
}
