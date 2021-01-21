package employee

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

// Employees is an aray of type Employee
type Employees struct {
	Data []*Employee `json:"data"`
}
