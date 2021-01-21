package project

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/maanijou/project-manager/database"
	"github.com/maanijou/project-manager/employee"
	"github.com/maanijou/project-manager/entity"
)

// insertProject is a database handler
func insertProject(project entity.Project) (int, error) {
	// NOTE: Better to use tx, err := db.Begin() as so many commits are involved
	// Right now just the participants with successful operations will be added
	owner, err := employee.GetEmployee(project.Owner.ID)
	if err != nil {
		log.Printf("Error in geting employee using uuid %v", project.Owner.ID)
		return 0, err
	}
	if string(owner.Role) != entity.ManagerRole.String() {
		log.Printf("Cannot set owner while the role is not Manager! %v", owner)
		return 0, errors.New("OwnerRoleProblem")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	var mu sync.Mutex
	var participantIDS []string

	for _, emp := range project.Participants {
		wg.Add(1)
		go func(ee entity.Employee, safe *sync.Mutex) {
			e, err := employee.GetEmployee(ee.ID)
			if err != nil {
				log.Printf("Error in geting employee using uuid %v", project.Owner.ID)
				wg.Done()
				return
			}
			if e.Department != owner.Department {
				log.Printf("Participant department is different from the owner %s != %s\n", e.Department, owner.Department)
				wg.Done()
				return
			}
			_, err = employee.InsertEmployee(&ee)
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") {
					mu.Lock()
					participantIDS = append(participantIDS, ee.ID)
					mu.Unlock()
				}
				wg.Done()
				return
			}
			mu.Lock()
			participantIDS = append(participantIDS, ee.ID)
			mu.Unlock()
			wg.Done()
		}(emp, &mu)
	}
	employee.InsertEmployee(&project.Owner)
	wg.Wait()
	ownerID, err := uuid.FromString(project.Owner.ID)
	if err != nil {
		return 0, err
	}
	stm, err := database.DbConn.Prepare(`INSERT INTO project
	(name, state, progress, owner)
	VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		log.Printf("Error inserting project %v\n", err)
		return 0, err
	}
	lastInsertID := 0
	err = stm.QueryRowContext(ctx, project.Name, project.State.Int(), project.Progress, ownerID).Scan(&lastInsertID)

	if err != nil {
		log.Printf("Error inserting project %v\n", err)
		return 0, err
	}
	for _, id := range participantIDS {
		stm, err := database.DbConn.Prepare(`INSERT INTO public.project_employee(
			project_id, employee_id)
			VALUES ($1, $2);`)
		if err != nil {
			log.Printf("Error inserting project %v\n", err)
			return 0, err
		}
		_, err = stm.ExecContext(ctx, lastInsertID, id)
		if err != nil {
			log.Printf("Error in inserting relation %d, %s", lastInsertID, id)
			return 0, err
		}
	}

	return lastInsertID, nil
}
