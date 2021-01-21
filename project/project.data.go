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
		stm, err := database.DbConn.Prepare(`INSERT INTO project_employee(
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

func getProjects(page, limit int) ([]entity.Project, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("Error in page or limit arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT id, name, state, progress, owner
	FROM public.project
	ORDER BY id ASC LIMIT $1 OFFSET $2
	`, limit, (page-1)*limit)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	projects := make([]entity.Project, 0)
	for results.Next() {
		var project entity.Project
		var state *int
		results.Scan(&project.ID,
			&project.Name,
			&state,
			&project.Progress,
			&project.Owner.ID)
		project.State = entity.StateName[*state]
		owner, err := employee.GetEmployee(project.Owner.ID)
		if err != nil {
			log.Printf("Error in geting employee using uuid %v", project.Owner.ID)
			return nil, err
		}
		project.Owner = *owner
		projects = append(projects, project)
	}
	// TODO Just ignore the participantsin for now. They can use /project/{id}
	// without preloading participants
	return projects, nil
}

func getProject(projectID int) (*entity.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(ctx, `SELECT id, name, state, progress, owner FROM project
	WHERE id = $1`, projectID)
	project := &entity.Project{}
	var state *int

	err := row.Scan(
		&project.ID,
		&project.Name,
		&state,
		&project.Progress,
		&project.Owner.ID,
	)

	if err != nil {
		log.Printf("Error in geting project %v", err)
		return nil, err
	}
	project.State = entity.StateName[*state]

	owner, err := employee.GetEmployee(project.Owner.ID)
	if err != nil {
		log.Printf("Error in geting employee using uuid %v", project.Owner.ID)
		return nil, err
	}
	project.Owner = *owner
	results, err := database.DbConn.QueryContext(ctx, `SELECT project_employee.employee_id FROM project 
	join project_employee ON project_employee.project_id = project.id
	WHERE project.id = $1`, projectID)
	defer results.Close()
	if err != nil {
		log.Println("Error in getProject", err)
		return nil, err
	}
	var employeeIDS []string
	for results.Next() {
		var id string
		results.Scan(&id)
		employeeIDS = append(employeeIDS, id)
	}
	for _, id := range employeeIDS {
		emp, err := employee.GetEmployee(id)
		if err != nil {
			log.Printf("Error in geting employee using uuid %v", id)
		}
		project.Participants = append(project.Participants, *emp)
	}
	return project, nil
}

func removeProject(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM project where id = $1 RETURNING id`, productID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// updateProject is used to just update the project, not its participants
func updateProject(project entity.Project) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if project.ID == 0 || *&project.ID == 0 {
		return -1, errors.New("product has invalid ID")
	}

	stm, err := database.DbConn.Prepare(`UPDATE project
	SET id=$1, name=$2, state=$3, progress=$4, owner=$5
	WHERE project.id = $6 RETURNING id;`)
	if err != nil {
		log.Printf("Error inserting project %v\n", err)
		return 0, err
	}
	lastInsertID := 0
	err = stm.QueryRowContext(ctx, project.ID,
		project.Name,
		project.State.Int(),
		project.Progress,
		project.Owner.ID,
		project.ID).Scan(&lastInsertID)

	if err != nil {
		log.Println(err.Error())
		return lastInsertID, err
	}
	return lastInsertID, err
}

func addParticipant(projectID int, participants []entity.Employee) (int, error) {
	project, err := getProject(projectID)
	if err != nil {
		return 0, err
	}
	owner := project.Owner
	wg := sync.WaitGroup{}
	var mu sync.Mutex
	var participantIDS []string

	for _, emp := range participants {
		log.Println(emp.ID)
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
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	count := 0
	for _, id := range participantIDS {
		stm, err := database.DbConn.Prepare(`INSERT INTO project_employee(
			project_id, employee_id)
			VALUES ($1, $2);`)
		if err != nil {
			log.Printf("Error inserting project %v\n", err)
			return 0, err
		}
		_, err = stm.ExecContext(ctx, project.ID, id)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				log.Printf("Error in inserting relation %d, %s %v", project.ID, id, err)
				return 0, err
			}
			continue
		}
		count = count + 1
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

func removeParticipant(projectID int, participantID string) error {
	_, err := getProject(projectID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stm, err := database.DbConn.Prepare(`DELETE FROM project_employee 
	WHERE project_employee.project_id = $1 and
	project_employee.employee_id = $2`)
	if err != nil {
		log.Printf("Error removing project %v\n", err)
		return err
	}
	_, err = stm.ExecContext(ctx, projectID, participantID)
	if err != nil {
		log.Printf("Error removing project %v\n", err)
		return err
	}
	return nil
}
