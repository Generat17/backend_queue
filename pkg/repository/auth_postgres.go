package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"server/types"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// CreateEmployee создает запись о новом сотруднике
func (r *AuthPostgres) CreateEmployee(username, password, firstName, secondName string, isAdmin bool) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password, first_name, second_name, is_admin, expires_at) values ($1, $2, $3, $4, $5, 0) RETURNING employee_id", employeeTable)

	row := r.db.QueryRow(query, username, password, firstName, secondName, isAdmin)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// GetEmployeeId получает ID сотрудника из БД по его логину и паролю
func (r *AuthPostgres) GetEmployeeId(username, password string) (types.Employee, error) {
	var employee types.Employee
	query := fmt.Sprintf("SELECT employee_id FROM %s WHERE username=$1 AND password=$2", employeeTable)
	err := r.db.Get(&employee, query, username, password)

	return employee, err
}

// GetStatusEmployee получает текущий статус сотрудника из БД
func (r *AuthPostgres) GetStatusEmployee(employeeId int) (int, error) {
	var employee types.Employee
	query := fmt.Sprintf("SELECT status FROM %s WHERE employee_id=$1", employeeTable)
	err := r.db.Get(&employee, query, employeeId)

	return employee.Status, err
}

// GetEmployeeById получает данные сотрудника из БД по его ID
func (r *AuthPostgres) GetEmployeeById(employeeId int) (types.Employee, error) {
	var employee types.Employee
	query := fmt.Sprintf("SELECT * FROM %s WHERE employee_id=$1", employeeTable)
	err := r.db.Get(&employee, query, employeeId)

	return employee, err
}

// GetEmployee получает данные сотрудника по его логину и паролю
func (r *AuthPostgres) GetEmployee(username, password string) (types.Employee, error) {
	var employee types.Employee
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 AND password=$2", employeeTable)
	err := r.db.Get(&employee, query, username, password)

	return employee, err
}

// SetSession записывает новое значение сессии для сотрудника
func (r *AuthPostgres) SetSession(refreshToken string, expiresAt int64, workstationId int, employeeId int) (sql.Result, error) {

	query := fmt.Sprintf("UPDATE %s SET refresh_token=$1, expires_at=$2, workstation_id=$3 WHERE employee_id=$4", employeeTable)
	res, err := r.db.Exec(query, refreshToken, expiresAt, workstationId, employeeId)

	return res, err
}

// GetSession получает текущее значение сессии сотрудника из БД
func (r *AuthPostgres) GetSession(employeeId int) (types.SessionInfo, error) {

	var sessionInfo types.SessionInfo
	query := fmt.Sprintf("SELECT refresh_token, expires_at, workstation_id FROM %s WHERE employee_id=$1", employeeTable)
	err := r.db.Get(&sessionInfo, query, employeeId)

	return sessionInfo, err
}

// ClearSession обнуляет значения сессии сотрудника по его ID
func (r *AuthPostgres) ClearSession(employeeId int) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET refresh_token=$1, expires_at=$2, workstation_id=$3, status=$4 WHERE employee_id=$5", employeeTable)
	res, err := r.db.Exec(query, "", 0, -1, 0, employeeId)

	return res, err
}
