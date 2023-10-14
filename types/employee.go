package types

type EmployeeResponsibilityItem struct {
	Id   int    `json:"employee_responsibility_id"`
	Name string `json:"employee_responsibility_name"`
}

// Employee структура, которая соответствует таблице Employee в БД
type Employee struct {
	EmployeeId         int                          `json:"employee_id" db:"employee_id"`
	Username           string                       `json:"username" db:"username"`
	Password           string                       `json:"password" db:"password"`
	FirstName          string                       `json:"first_name" db:"first_name"`
	SecondName         string                       `json:"second_name" db:"second_name"`
	SessionState       bool                         `json:"session_state" db:"session_state"`
	Status             int                          `json:"status" db:"status"`
	RefreshToken       string                       `json:"refresh_token" db:"refresh_token"`
	ExpiresAt          int                          `json:"expires_at" db:"expires_at"`
	WorkstationId      int                          `json:"workstation_id" db:"workstation_id"`
	IsAdmin            bool                         `json:"is_admin" db:"is_admin"`
	ResponsibilityList []EmployeeResponsibilityItem `json:"responsibility_list"`
}

type EmployeeStatus struct {
	EmployeeId int `json:"employee_id" db:"employee_id"`
	Status     int `json:"status" db:"status"`
}

type GetNewClientResponse struct {
	NumberTicket   int    `json:"number_ticket"`
	ServiceTicket  string `json:"service_ticket"`
	EmployeeStatus int    `json:"employee_status"`
	NumberQueue    int    `json:"number_queue"`
}

type GetEmployeeListsResponse struct {
	Data []Employee `json:"data"`
}

type ConfirmClientResponse struct {
	NumberQueue int `json:"number_queue"`
}

type EmployeeStatusResponse struct {
	EmployeeStatus int `json:"employee_status"`
}

type EmployeeResponsibility struct {
	EmployeeId         int    `json:"employee_id"  db:"employee_id"`
	ResponsibilityId   int    `json:"responsibility_id"  db:"responsibility_id"`
	ResponsibilityName string `json:"responsibility_name"  db:"responsibility_name"`
}
