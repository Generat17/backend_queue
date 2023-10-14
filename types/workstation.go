package types

import "database/sql"

type WorkstationResponsibilityItem struct {
	Id   int    `json:"workstation_responsibility_id"`
	Name string `json:"workstation_responsibility_name"`
}

type Workstation struct {
	WorkstationId      int                             `json:"workstation_id"  db:"workstation_id"`
	WorkstationName    string                          `json:"workstation_name"  db:"workstation_name"`
	EmployeeId         int                             `json:"employee_id"  db:"employee_id"`
	IsBusy             bool                            `json:"is_busy"  db:"is_busy"`
	ResponsibilityList []WorkstationResponsibilityItem `json:"responsibility_list"`
}

type WorkstationResponsibility struct {
	WorkstationId      int    `json:"workstation_id"  db:"workstation_id"`
	ResponsibilityId   int    `json:"responsibility_id"  db:"responsibility_id"`
	ResponsibilityName string `json:"responsibility_name"  db:"responsibility_name"`
}

type ResponseWorkstation struct {
	Response sql.Result `json:"response"`
}
