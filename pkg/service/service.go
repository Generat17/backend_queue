package service

import (
	"database/sql"
	"server/pkg/repository"
	"server/types"
)

type Authorization interface {
	CreateEmployee(username, password, firstName, secondName string, isAdmin bool) (int, error)
	GenerateTokenWorkstation(username, password string, workstation int) (string, error)
	ParseTokenWorkstation(token string) (types.ParseTokenWorkstationResponse, error)
	GenerateRefreshToken() (string, error)
	UpdateTokenWorkstation(employeeId, workstationId int, refreshToken string) (string, error)
	SetSession(refreshToken string, workstationId int, employeeId int) (bool, error)
	LogOut(employeeId int) (bool, error)
	GetEmployee(username, password string) (types.Employee, error)
	GetEmployeeById(employeeId int) (types.Employee, error)
	GetStatusEmployee(employeeId int) (int, error)
}

type Employee interface {
	GetEmployeeList() ([]types.Employee, error)
	UpdateEmployeeResponsibility(employeeId int, responsibilityIdList []int) ([]types.Employee, error)
	UpdateEmployee(employeeId int, username, firstName, secondName string, isAdmin bool) (sql.Result, error)
	RemoveEmployee(employeeId int) (sql.Result, error)
	GetEmployeeStatusList() ([]types.EmployeeStatus, error)
	GetEmployeeStatus(workstationId int) (types.EmployeeStatus, error)
}

type Queue interface {
	GetQueueList() ([]types.QueueItem, error)
	GetQueueAdminList() ([]types.QueueItem, error)
	GetQueueItemStatus(workstationId int) (types.QueueItem, error)
	AddQueueItem(service string) (int, error)
	GetNewClient(employeeId, workstationId int) (types.GetNewClientResponse, error)
	ConfirmClient(numberQueue, employeeId int) (int, error)
	EndClient(numberQueue, employeeId int) (int, error)
	SetEmployeeStatus(statusCode, employeeId int) (bool, error)
	ClearQueue() (bool, error)
	ClearLog() (bool, error)
	UpdateQuality(quality, client int) (bool, error)
	GetClientsLog() ([]types.LogItem, error)
	NotCome(numberQueue, employeeId int) (int, error)
	SendEmail(to []string, subject, body string) (bool, error)
	CheckLongWait() (bool, error)
	GetTiming() ([]types.TimingItem, error)
	GetEmail() ([]types.EmailItem, error)
	UpdateTiming(id int, seconds int, name string) (sql.Result, error)
	RemoveTiming(id int) (sql.Result, error)
	AddTiming(seconds int, name string) (sql.Result, error)
	RemoveEmail(id int) (sql.Result, error)
	AddEmail(timing int, email string) (sql.Result, error)
	ActiveTiming(id int) (sql.Result, error)
	RestartIdentity() (sql.Result, error)
}

type Responsibility interface {
	GetResponsibilityList() ([]types.Responsibility, error)
	UpdateResponsibility(responsibilityId int, responsibilityName string, responsibilityPriority int) (sql.Result, error)
	RemoveResponsibility(responsibilityId int) (sql.Result, error)
	AddResponsibility(responsibilityName string) (sql.Result, error)
}

type Workstation interface {
	GetWorkstationList() ([]types.Workstation, error)
	GetWorkstation(workstationId int) (types.Workstation, error)
	AddWorkstation(workstationName string) (sql.Result, error)
	UpdateWorkstation(workstationId int, workstationName string) (sql.Result, error)
	RemoveWorkstation(workstationId int) (sql.Result, error)
	UpdateWorkstationResponsibility(workstationId int, responsibilityId []int) ([]types.Workstation, error)
	UpdateWorkstationStatus(workstationId int, workstationStatus bool, employeeId int) (sql.Result, error)
}

type Service struct {
	Employee
	Queue
	Responsibility
	Authorization
	Workstation
}

func NewService(repos *repository.Repository, checkUpdate *[1]int) *Service {
	return &Service{
		Authorization:  NewAuthService(repos.Authorization),
		Employee:       NewEmployeeService(repos.Employee),
		Queue:          NewQueueService(repos.Queue, checkUpdate),
		Responsibility: NewResponsibilityService(repos.Responsibility),
		Workstation:    NewWorkstationService(repos.Workstation),
	}
}
