package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"server/types"
)

type EmployeePostgres struct {
	db *sqlx.DB
}

func NewEmployeePostgres(db *sqlx.DB) *EmployeePostgres {
	return &EmployeePostgres{db: db}
}

// GetEmployeeList получает список пользователей из БД
func (r *EmployeePostgres) GetEmployeeList() ([]types.Employee, error) {
	var employee []types.Employee
	query := fmt.Sprintf("SELECT * FROM %s", employeeTable)
	err := r.db.Select(&employee, query)

	return employee, err
}

// GetEmployeeStatusList получает список сотрудников и их текущий статус из бд
func (r *EmployeePostgres) GetEmployeeStatusList() ([]types.EmployeeStatus, error) {
	var employeeStatus []types.EmployeeStatus
	query := fmt.Sprintf("SELECT employee_id, status FROM %s", employeeTable)
	err := r.db.Select(&employeeStatus, query)

	return employeeStatus, err
}

// GetEmployeeStatus получает статус сотрудника и его текущий статус из бд
func (r *EmployeePostgres) GetEmployeeStatus(workstationId int) ([]types.EmployeeStatus, error) {
	var employeeStatus []types.EmployeeStatus

	query := fmt.Sprintf("SELECT employee_id, status FROM %s WHERE workstation_id = $1", employeeTable)
	err := r.db.Select(&employeeStatus, query, workstationId)

	return employeeStatus, err
}

// GetEmployeeResponsibilityList получает список рабочих станций и обязанностей для них из БД
func (r *EmployeePostgres) GetEmployeeResponsibilityList() ([]types.EmployeeResponsibility, error) {
	var employee []types.EmployeeResponsibility
	query := fmt.Sprintf("SELECT em.employee_id, em.responsibility_id, r.responsibility_name from employee_responsibility as em join responsibility as r on em.responsibility_id = r.responsibility_id")
	err := r.db.Select(&employee, query)

	return employee, err
}

// AddEmployeeResponsibility добавляет запись рабочая станция - обязанность в БД
func (r *EmployeePostgres) AddEmployeeResponsibility(employeeId int, responsibilityId int) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (employee_id, responsibility_id) VALUES ($1, $2)", employeeResponsibilityTable)
	res, err := r.db.Exec(query, employeeId, responsibilityId)

	return res, err
}

// RemoveEmployeeResponsibility удаляет запись рабочая станция - обязанность в БД
func (r *EmployeePostgres) RemoveEmployeeResponsibility(employeeId int, responsibilityId int) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE employee_id=$1 AND responsibility_id=$2", employeeResponsibilityTable)
	res, err := r.db.Exec(query, employeeId, responsibilityId)

	return res, err
}

// GetEmployeeResponsibilityById получает список рабочих станций и обязанностей для них из БД
func (r *EmployeePostgres) GetEmployeeResponsibilityById(employeeId int) ([]int, error) {
	var responsibilityIdList []int
	query := fmt.Sprintf("SELECT responsibility_id FROM %s WHERE employee_id = $1", employeeResponsibilityTable)
	err := r.db.Select(&responsibilityIdList, query, employeeId)

	return responsibilityIdList, err
}

// UpdateEmployee обновляет запись о сотруднике в БД по id
func (r *EmployeePostgres) UpdateEmployee(employeeId int, username, firstName, secondName string, isAdmin bool) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET username=$1, first_name=$2, second_name=$3, is_admin=$4 WHERE employee_id=$5", employeeTable)
	res, err := r.db.Exec(query, username, firstName, secondName, isAdmin, employeeId)

	return res, err
}

// RemoveEmployee удаляет запись о сотруднике из БД по id
func (r *EmployeePostgres) RemoveEmployee(employeeId int) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE employee_id=$1", employeeTable)
	res, err := r.db.Exec(query, employeeId)

	return res, err
}
