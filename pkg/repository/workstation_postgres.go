package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"server/types"
)

type WorkstationPostgres struct {
	db *sqlx.DB
}

func NewWorkstationPostgres(db *sqlx.DB) *WorkstationPostgres {
	return &WorkstationPostgres{db: db}
}

// GetWorkstationList получает список рабочих станций из БД
func (r *WorkstationPostgres) GetWorkstationList() ([]types.Workstation, error) {
	var workstation []types.Workstation
	query := fmt.Sprintf("SELECT * FROM %s", workstationTable)
	err := r.db.Select(&workstation, query)

	return workstation, err
}

// GetWorkstationResponsibilityList получает список рабочих станций и обязанностей для них из БД
func (r *WorkstationPostgres) GetWorkstationResponsibilityList() ([]types.WorkstationResponsibility, error) {
	var workstation []types.WorkstationResponsibility
	query := fmt.Sprintf("SELECT wr.workstation_id, wr.responsibility_id, r.responsibility_name from workstation_responsibility as wr join responsibility as r on wr.responsibility_id = r.responsibility_id")
	err := r.db.Select(&workstation, query)

	return workstation, err
}

// GetWorkstation получает данные о рабочем месте из БД по его ID
func (r *WorkstationPostgres) GetWorkstation(workstationId int) (types.Workstation, error) {
	var workstation types.Workstation
	query := fmt.Sprintf("SELECT workstation_id, workstation_name, employee_id, is_busy FROM %s WHERE workstation_id=$1", workstationTable)
	err := r.db.Get(&workstation, query, workstationId)

	return workstation, err
}

// UpdateWorkstation обновляет запись о рабочей станции в БД по id
func (r *WorkstationPostgres) UpdateWorkstation(workstationId int, workstationName string) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET workstation_name=$1 WHERE workstation_id=$2", workstationTable)
	res, err := r.db.Exec(query, workstationName, workstationId)

	return res, err
}

// RemoveWorkstation удаляет запись о рабочей станции в БД по id
func (r *WorkstationPostgres) RemoveWorkstation(workstationId int) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE workstation_id=$1", workstationTable)
	res, err := r.db.Exec(query, workstationId)

	return res, err
}

// AddWorkstation добавляет запись о рабочей станции в БД
func (r *WorkstationPostgres) AddWorkstation(workstationName string) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (workstation_name) VALUES ($1)", workstationTable)
	res, err := r.db.Exec(query, workstationName)

	return res, err
}

// GetWorkstationResponsibilityById получает список рабочих станций и обязанностей для них из БД
func (r *WorkstationPostgres) GetWorkstationResponsibilityById(workstationId int) ([]int, error) {
	var responsibilityIdList []int
	query := fmt.Sprintf("SELECT responsibility_id FROM %s WHERE workstation_id = $1", workstationResponsibilityTable)
	err := r.db.Select(&responsibilityIdList, query, workstationId)

	return responsibilityIdList, err
}

// AddWorkstationResponsibility добавляет запись рабочая станция - обязанность в БД
func (r *WorkstationPostgres) AddWorkstationResponsibility(workstationId int, responsibilityId int) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (workstation_id, responsibility_id) VALUES ($1, $2)", workstationResponsibilityTable)
	res, err := r.db.Exec(query, workstationId, responsibilityId)

	return res, err
}

// RemoveWorkstationResponsibility удаляет запись рабочая станция - обязанность в БД
func (r *WorkstationPostgres) RemoveWorkstationResponsibility(workstationId int, responsibilityId int) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE workstation_id=$1 AND responsibility_id=$2", workstationResponsibilityTable)
	res, err := r.db.Exec(query, workstationId, responsibilityId)

	return res, err
}

// UpdateWorkstationStatus обновляет запись о рабочей станции в БД по id
func (r *WorkstationPostgres) UpdateWorkstationStatus(workstationId int, workstationStatus bool, employeeId int) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET is_busy=$1, employee_id=$2 WHERE workstation_id=$3", workstationTable)
	res, err := r.db.Exec(query, workstationStatus, employeeId, workstationId)

	return res, err
}
