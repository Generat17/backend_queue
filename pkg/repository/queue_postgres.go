package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"server/types"
	"time"
)

type QueuePostgres struct {
	db *sqlx.DB
}

func NewQueuePostgres(db *sqlx.DB) *QueuePostgres {
	return &QueuePostgres{db: db}
}

// GetQueue получает текущее состояние очереди из БД
func (r *QueuePostgres) GetQueue() ([]types.QueueItem, error) {

	var queue []types.QueueItem
	query := fmt.Sprintf("SELECT * FROM %s", queueTable)
	err := r.db.Select(&queue, query)

	return queue, err
}

// GetEmail получает текущее состояние очереди из БД
func (r *QueuePostgres) GetEmail() ([]types.EmailItem, error) {

	var email []types.EmailItem
	query := fmt.Sprintf("SELECT * FROM %s", emailTable)
	err := r.db.Select(&email, query)

	return email, err
}

// GetTiming получает текущее состояние очереди из БД
func (r *QueuePostgres) GetTiming() ([]types.TimingItem, error) {

	var timing []types.TimingItem
	query := fmt.Sprintf("SELECT * FROM %s", timingTable)
	err := r.db.Select(&timing, query)

	return timing, err
}

// GetQueueItemStatus получает текущее состояние элемента в очереди с соответствующим id workstation
func (r *QueuePostgres) GetQueueItemStatus(workstationId int) ([]types.QueueItem, error) {

	var queue []types.QueueItem
	query := fmt.Sprintf("SELECT * FROM %s WHERE workstation_number = $1", queueTable)
	err := r.db.Select(&queue, query, workstationId)

	return queue, err
}

// AddItemQueue добавляет нового клиента в очередь
func (r *QueuePostgres) AddItemQueue(service string) (int, error) {

	// Получение текущего времени
	time := time.Now().Unix()

	var itemId int

	var responsibility []types.Responsibility
	query := fmt.Sprintf("SELECT * FROM %s WHERE responsibility_name=$1", responsibilityTable)
	errSelect := r.db.Select(&responsibility, query, service)
	if errSelect != nil {
		logrus.Print(errSelect)
	}

	err := r.db.QueryRowx(`INSERT INTO queue (priority, service_type, workstation_number, status, start_time, call_time, service_time) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_client`, responsibility[0].ResponsibilityPriority, service, -1, 1, time, -1, -1).Scan(&itemId)
	if err != nil {
		logrus.Print(err)
		return -1, err
	}

	return itemId, err
}

// CallClientTime записывает время, когда оператор вызвал
func (r *QueuePostgres) CallClientTime(numberQueue, workstationId int) (sql.Result, error) {

	// Получение текущего времени
	time := time.Now().Unix()

	query := fmt.Sprintf("UPDATE %s SET call_time=$1, workstation_number=$2 WHERE id_client=$3", queueTable)
	res, err := r.db.Exec(query, time, workstationId, numberQueue)

	return res, err
}

// ConfirmClientTime записывает время, когда клиент подошел к оператору
func (r *QueuePostgres) ConfirmClientTime(numberQueue int) (sql.Result, error) {

	// Получение текущего времени
	time := time.Now().Unix()

	query := fmt.Sprintf("UPDATE %s SET service_time=$1 WHERE id_client=$2", queueTable)
	res, err := r.db.Exec(query, time, numberQueue)

	return res, err
}

// NotCome клиент не пришел
func (r *QueuePostgres) NotCome(numberQueue int) (sql.Result, error) {
	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE id_client=$1", queueTable)
	res, err := r.db.Exec(queryDelete, numberQueue)

	return res, err
}

// EndClientTime записывает время, когда обслуживание клиента завершилось
func (r *QueuePostgres) EndClientTime(numberQueue int) (sql.Result, error) {

	var client []types.QueueItem
	var workstation []types.Workstation
	var employee []types.Employee

	querySelect1 := fmt.Sprintf("SELECT * FROM %s WHERE id_client = $1", queueTable)
	errSelect1 := r.db.Select(&client, querySelect1, numberQueue)
	if errSelect1 != nil {
		logrus.Print(errSelect1)
	}

	querySelect2 := fmt.Sprintf("SELECT * FROM %s WHERE workstation_id = $1", workstationTable)
	errSelect2 := r.db.Select(&workstation, querySelect2, client[0].WorkstationNumber)
	if errSelect2 != nil {
		logrus.Print(errSelect2)
	}

	querySelect3 := fmt.Sprintf("SELECT * FROM %s WHERE employee_id = $1", employeeTable)
	errSelect3 := r.db.Select(&employee, querySelect3, workstation[0].EmployeeId)
	if errSelect3 != nil {
		logrus.Print(errSelect3)
	}

	time := time.Now().Unix()

	errInsert := r.db.QueryRowx(`INSERT INTO clients_log (ticket_number, priority, service_type, workstation_number, status, start_time, call_time, service_time, end_time, quality, employee_id, employee_first_name, employee_second_name, workstation_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, client[0].Id, client[0].Priority, client[0].ServiceType, client[0].WorkstationNumber, client[0].Status, client[0].StartTime, client[0].CallTime, client[0].ServiceTime, time, client[0].Quality, employee[0].EmployeeId, employee[0].FirstName, employee[0].SecondName, workstation[0].WorkstationName)
	if errInsert != nil {
		logrus.Print(errInsert)
	}

	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE id_client=$1", queueTable)
	res, err := r.db.Exec(queryDelete, numberQueue)

	return res, err
}

// SetStatusEmployee обновляет статус сотрудника в БД
func (r *QueuePostgres) SetStatusEmployee(statusCode int, employeeId int) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET status=$1 WHERE employee_id=$2", employeeTable)
	res, err := r.db.Exec(query, statusCode, employeeId)

	return res, err
}

// GetResponsibilityByEmployeeId получает список обязанностей сотрудника по его ID
func (r *QueuePostgres) GetResponsibilityByEmployeeId(employeeId int) ([]types.Responsibility, error) {

	var responsibility []types.Responsibility
	query := fmt.Sprintf("SELECT responsibility.responsibility_id, responsibility_name, responsibility_priority FROM employee_responsibility LEFT JOIN responsibility ON employee_responsibility.responsibility_id = responsibility.responsibility_id WHERE employee_responsibility.employee_id = $1")
	err := r.db.Select(&responsibility, query, employeeId)

	return responsibility, err
}

// GetResponsibilityByWorkstationId получает функции (обязанности) рабочего места по его ID
func (r *QueuePostgres) GetResponsibilityByWorkstationId(workstationId int) ([]types.Responsibility, error) {

	var responsibility []types.Responsibility
	query := fmt.Sprintf("SELECT responsibility.responsibility_id, responsibility_name, responsibility_priority FROM workstation_responsibility LEFT JOIN responsibility ON workstation_responsibility.responsibility_id = responsibility.responsibility_id WHERE workstation_responsibility.workstation_id = $1")
	err := r.db.Select(&responsibility, query, workstationId)

	return responsibility, err
}

// ClearQueue получает функции (обязанности) рабочего места по его ID
func (r *QueuePostgres) ClearQueue() (sql.Result, error) {

	query := fmt.Sprintf("DELETE FROM %s", queueTable)
	res, err := r.db.Exec(query)

	return res, err
}

// ClearLog получает функции (обязанности) рабочего места по его ID
func (r *QueuePostgres) ClearLog() (sql.Result, error) {

	query := fmt.Sprintf("DELETE FROM %s", clientsLog)
	res, err := r.db.Exec(query)

	return res, err
}

// UpdateQuality записывает время, когда клиент подошел к оператору
func (r *QueuePostgres) UpdateQuality(quality, client int) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET quality=$1 WHERE id_client=$2", queueTable)
	res, err := r.db.Exec(query, quality, client)

	return res, err
}

// GetLogs записывает время, когда клиент подошел к оператору
func (r *QueuePostgres) GetLogs() ([]types.LogItem, error) {
	var logs []types.LogItem

	querySelect := fmt.Sprintf("SELECT * FROM %s", clientsLog)
	err := r.db.Select(&logs, querySelect)
	if err != nil {
		logrus.Print(err)
	}

	return logs, err
}

// UpdateTiming обновляет запись об Обязанности в БД по id
func (r *QueuePostgres) UpdateTiming(id int, seconds int, name string) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET seconds=$1, name=$2 WHERE id=$3", timingTable)
	res, err := r.db.Exec(query, seconds, name, id)

	return res, err
}

// RemoveTiming удаляет запись об Обязанности в БД по id
func (r *QueuePostgres) RemoveTiming(id int) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", timingTable)
	res, err := r.db.Exec(query, id)

	return res, err
}

// AddTiming добавляет запись об Обязанности в БД
func (r *QueuePostgres) AddTiming(seconds int, name string) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (seconds, name) VALUES ($1,$2)", timingTable)
	res, err := r.db.Exec(query, seconds, name)

	return res, err
}

// RemoveEmail удаляет запись об Обязанности в БД по id
func (r *QueuePostgres) RemoveEmail(id int) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", emailTable)
	res, err := r.db.Exec(query, id)

	return res, err
}

// AddEmail добавляет запись об Обязанности в БД
func (r *QueuePostgres) AddEmail(timing int, email string) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (timing, email) VALUES ($1, $2)", emailTable)
	res, err := r.db.Exec(query, timing, email)

	return res, err
}

// ActiveTiming добавляет запись об Обязанности в БД
func (r *QueuePostgres) ActiveTiming(id int) (sql.Result, error) {
	var timing []types.TimingItem
	querySelect := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", timingTable)
	err := r.db.Select(&timing, querySelect, id)

	queryUpdate := fmt.Sprintf("UPDATE %s SET active=$1 WHERE id=$2", timingTable)
	res, err := r.db.Exec(queryUpdate, !timing[0].Active, id)

	return res, err
}

// SetStatusQueueItem обновляет статус сотрудника в БД
func (r *QueuePostgres) SetStatusQueueItem(statusCode int, clientId int) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET status=$1 WHERE id_client=$2", queueTable)
	res, err := r.db.Exec(query, statusCode, clientId)

	return res, err
}

// RestartIdentity обнуляет счетчик авто инкремента
func (r *QueuePostgres) RestartIdentity() (sql.Result, error) {
	queryDelete := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", queueTable)
	res, err := r.db.Exec(queryDelete)

	return res, err
}
