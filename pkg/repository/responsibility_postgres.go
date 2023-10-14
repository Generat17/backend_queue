package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"server/types"
)

type ResponsibilityPostgres struct {
	db *sqlx.DB
}

func NewResponsibilityPostgres(db *sqlx.DB) *ResponsibilityPostgres {
	return &ResponsibilityPostgres{db: db}
}

// GetResponsibilityList получает список всех обязанностей из БД
func (r *ResponsibilityPostgres) GetResponsibilityList() ([]types.Responsibility, error) {
	var responsibility []types.Responsibility
	query := fmt.Sprintf("SELECT * FROM %s", responsibilityTable)
	err := r.db.Select(&responsibility, query)

	return responsibility, err
}

// UpdateResponsibility обновляет запись об Обязанности в БД по id
func (r *ResponsibilityPostgres) UpdateResponsibility(responsibilityId int, responsibilityName string, responsibilityPriority int) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET responsibility_name=$1, responsibility_priority=$2 WHERE responsibility_id=$3", responsibilityTable)
	res, err := r.db.Exec(query, responsibilityName, responsibilityPriority, responsibilityId)

	return res, err
}

// RemoveResponsibility удаляет запись об Обязанности в БД по id
func (r *ResponsibilityPostgres) RemoveResponsibility(responsibilityId int) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE responsibility_id=$1", responsibilityTable)
	res, err := r.db.Exec(query, responsibilityId)

	return res, err
}

// AddResponsibility добавляет запись об Обязанности в БД
func (r *ResponsibilityPostgres) AddResponsibility(responsibilityName string) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (responsibility_name, responsibility_priority) VALUES ($1, 1)", responsibilityTable)
	res, err := r.db.Exec(query, responsibilityName)

	return res, err
}
