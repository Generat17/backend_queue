package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"server/types"
)

type Authorization interface {
	CreateEmployee(username, password, firstName, secondName string, isAdmin bool) (int, error)
	GetEmployeeId(username, password string) (types.Employee, error)
	GetEmployee(username, password string) (types.Employee, error)
	GetEmployeeById(employeeId int) (types.Employee, error)
	SetSession(refreshToken string, expiresAt int64, workstationId int, employeeId int) (sql.Result, error)
	GetSession(employeeId int) (types.SessionInfo, error)
	ClearSession(employeeId int) (sql.Result, error)
	GetStatusEmployee(employeeId int) (int, error)
}

type Employee interface {
	GetEmployeeList() ([]types.Employee, error)
	GetEmployeeStatusList() ([]types.EmployeeStatus, error)
	GetEmployeeStatus(workstationId int) ([]types.EmployeeStatus, error)
	GetEmployeeResponsibilityList() ([]types.EmployeeResponsibility, error)
	AddEmployeeResponsibility(employeeId int, responsibilityId int) (sql.Result, error)
	RemoveEmployeeResponsibility(employeeId int, responsibilityId int) (sql.Result, error)
	GetEmployeeResponsibilityById(employeeId int) ([]int, error)
	UpdateEmployee(employeeId int, username, firstName, secondName string, isAdmin bool) (sql.Result, error)
	RemoveEmployee(employeeId int) (sql.Result, error)
}

type Responsibility interface {
	GetResponsibilityList() ([]types.Responsibility, error)
	RemoveResponsibility(responsibilityId int) (sql.Result, error)
	UpdateResponsibility(responsibilityId int, responsibilityName string, responsibilityPriority int) (sql.Result, error)
	AddResponsibility(responsibilityName string) (sql.Result, error)
}

type Queue interface {
	GetResponsibilityByEmployeeId(employeeId int) ([]types.Responsibility, error)
	GetResponsibilityByWorkstationId(workstationId int) ([]types.Responsibility, error)
	SetStatusEmployee(statusCode int, employeeId int) (sql.Result, error)
	GetQueue() ([]types.QueueItem, error)
	GetQueueItemStatus(workstationId int) ([]types.QueueItem, error)
	AddItemQueue(service string) (int, error)
	CallClientTime(numberQueue, workstationId int) (sql.Result, error)
	ConfirmClientTime(numberQueue int) (sql.Result, error)
	EndClientTime(numberQueue int) (sql.Result, error)
	ClearQueue() (sql.Result, error)
	ClearLog() (sql.Result, error)
	UpdateQuality(quality, client int) (sql.Result, error)
	GetLogs() ([]types.LogItem, error)
	NotCome(numberQueue int) (sql.Result, error)
	GetEmail() ([]types.EmailItem, error)
	GetTiming() ([]types.TimingItem, error)
	UpdateTiming(id int, seconds int, name string) (sql.Result, error)
	RemoveTiming(id int) (sql.Result, error)
	AddTiming(seconds int, name string) (sql.Result, error)
	RemoveEmail(id int) (sql.Result, error)
	AddEmail(timing int, email string) (sql.Result, error)
	ActiveTiming(id int) (sql.Result, error)
	SetStatusQueueItem(statusCode int, clientId int) (sql.Result, error)
	RestartIdentity() (sql.Result, error)
}

type Workstation interface {
	GetWorkstationList() ([]types.Workstation, error)
	GetWorkstation(workstationId int) (types.Workstation, error)
	RemoveWorkstation(workstationId int) (sql.Result, error)
	UpdateWorkstation(workstationId int, workstationName string) (sql.Result, error)
	AddWorkstation(workstationName string) (sql.Result, error)
	GetWorkstationResponsibilityList() ([]types.WorkstationResponsibility, error)
	AddWorkstationResponsibility(workstationId int, responsibilityId int) (sql.Result, error)
	RemoveWorkstationResponsibility(workstationId int, responsibilityId int) (sql.Result, error)
	GetWorkstationResponsibilityById(workstationId int) ([]int, error)
	UpdateWorkstationStatus(workstationId int, workstationStatus bool, employeeId int) (sql.Result, error)
}

type Repository struct {
	Employee
	Responsibility
	Authorization
	Queue
	Workstation
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:  NewAuthPostgres(db),
		Employee:       NewEmployeePostgres(db),
		Responsibility: NewResponsibilityPostgres(db),
		Queue:          NewQueuePostgres(db),
		Workstation:    NewWorkstationPostgres(db),
	}
}
