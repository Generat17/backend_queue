package types

type QueueItem struct {
	Id                int    `json:"id_client" db:"id_client"`
	Priority          int    `json:"priority" db:"priority"`
	ServiceType       string `json:"service_type" db:"service_type"`
	WorkstationNumber int    `json:"workstation_number" db:"workstation_number"`
	Status            int    `json:"status" db:"status"`
	StartTime         int    `json:"start_time" db:"start_time"`
	CallTime          int    `json:"call_time" db:"call_time"`
	ServiceTime       int    `json:"service_time" db:"service_time"`
	Quality           int    `json:"quality" db:"quality"`
}

type QueueItemNumber struct {
	Ticket int `json:"TicketID"`
}

type LogItem struct {
	Id                 int    `json:"id" db:"id"`
	TicketNumber       int    `json:"ticket_number" db:"ticket_number"`
	Priority           int    `json:"priority" db:"priority"`
	ServiceType        string `json:"service_type" db:"service_type"`
	WorkstationNumber  int    `json:"workstation_number" db:"workstation_number"`
	Status             int    `json:"status" db:"status"`
	StartTime          int    `json:"start_time" db:"start_time"`
	CallTime           int    `json:"call_time" db:"call_time"`
	ServiceTime        int    `json:"service_time" db:"service_time"`
	EndTime            int    `json:"end_time" db:"end_time"`
	Quality            int    `json:"quality" db:"quality"`
	EmployeeId         int    `json:"employee_id" db:"employee_id"`
	EmployeeFirstName  string `json:"employee_first_name" db:"employee_first_name"`
	EmployeeSecondName string `json:"employee_second_name" db:"employee_second_name"`
	WorkstationName    string `json:"workstation_name" db:"workstation_name"`
}

type EmailItem struct {
	Id     int    `json:"id" db:"id"`
	Timing int    `json:"timing" db:"timing"`
	Email  string `json:"email" db:"email"`
}

type TimingItem struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Seconds int    `json:"seconds" db:"seconds"`
	Active  bool   `json:"active" db:"active"`
}
