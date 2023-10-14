package service

import (
	"database/sql"
	"server/pkg/repository"
	"server/types"
)

type ResponsibilityService struct {
	repo repository.Responsibility
}

func NewResponsibilityService(repo repository.Responsibility) *ResponsibilityService {
	return &ResponsibilityService{repo: repo}
}

// GetResponsibilityList получает список обязанностей
func (s *ResponsibilityService) GetResponsibilityList() ([]types.Responsibility, error) {
	return s.repo.GetResponsibilityList()
}

// UpdateResponsibility обновляет обязанность
func (s *ResponsibilityService) UpdateResponsibility(responsibilityId int, responsibilityName string, responsibilityPriority int) (sql.Result, error) {
	return s.repo.UpdateResponsibility(responsibilityId, responsibilityName, responsibilityPriority)
}

// RemoveResponsibility удаляет обязанность
func (s *ResponsibilityService) RemoveResponsibility(responsibilityId int) (sql.Result, error) {
	return s.repo.RemoveResponsibility(responsibilityId)
}

// AddResponsibility создает обязанность
func (s *ResponsibilityService) AddResponsibility(responsibilityName string) (sql.Result, error) {
	return s.repo.AddResponsibility(responsibilityName)
}
