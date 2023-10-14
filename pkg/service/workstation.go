package service

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"server/pkg/repository"
	"server/types"
)

type WorkstationService struct {
	repo repository.Workstation
}

func NewWorkstationService(repo repository.Workstation) *WorkstationService {
	return &WorkstationService{repo: repo}
}

// Contains указывает, содержится ли x в a. Вспомогательная функция
func Contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// GetWorkstationList получает список рабочих мест
func (s *WorkstationService) GetWorkstationList() ([]types.Workstation, error) {

	listWorkstationResponsibility, err := s.repo.GetWorkstationResponsibilityList()
	if err != nil {
		return nil, err
	}

	listWorkstation, err := s.repo.GetWorkstationList()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(listWorkstation); i++ {

		for j := 0; j < len(listWorkstationResponsibility); j++ {
			if listWorkstationResponsibility[j].WorkstationId == listWorkstation[i].WorkstationId {
				listWorkstation[i].ResponsibilityList = append(listWorkstation[i].ResponsibilityList, types.WorkstationResponsibilityItem{Id: listWorkstationResponsibility[j].ResponsibilityId, Name: listWorkstationResponsibility[j].ResponsibilityName})
			}
		}
	}

	return listWorkstation, nil
}

// GetWorkstation получает данные о рабочем месте по его ID
func (s *WorkstationService) GetWorkstation(workstationId int) (types.Workstation, error) {
	return s.repo.GetWorkstation(workstationId)
}

// UpdateWorkstation обновляет рабочее место
func (s *WorkstationService) UpdateWorkstation(workstationId int, workstationName string) (sql.Result, error) {
	return s.repo.UpdateWorkstation(workstationId, workstationName)
}

// RemoveWorkstation удаляет рабочее место
func (s *WorkstationService) RemoveWorkstation(workstationId int) (sql.Result, error) {
	return s.repo.RemoveWorkstation(workstationId)
}

// AddWorkstation создает рабочее место
func (s *WorkstationService) AddWorkstation(workstationName string) (sql.Result, error) {
	return s.repo.AddWorkstation(workstationName)
}

// UpdateWorkstationResponsibility обновляет связь рабочее место - обязанность
func (s *WorkstationService) UpdateWorkstationResponsibility(workstationId int, responsibilityIdList []int) ([]types.Workstation, error) {
	currentWorkstationResponsibilityIdList, _ := s.repo.GetWorkstationResponsibilityById(workstationId) // получаем текущие связи

	var generalResponsibilityIdList = []int{} // в этом массиве храним общие функции рабочих мест

	// находим общие функции рабочих мест и записываем их в массив
	for i := 0; i < len(currentWorkstationResponsibilityIdList); i++ {
		for j := 0; j < len(responsibilityIdList); j++ {
			if currentWorkstationResponsibilityIdList[i] == responsibilityIdList[j] {
				generalResponsibilityIdList = append(generalResponsibilityIdList, currentWorkstationResponsibilityIdList[i])
			}
		}
	}

	// Удалим лишние элементы
	for i := 0; i < len(currentWorkstationResponsibilityIdList); i++ {
		if !Contains(generalResponsibilityIdList, currentWorkstationResponsibilityIdList[i]) {
			s.repo.RemoveWorkstationResponsibility(workstationId, currentWorkstationResponsibilityIdList[i])
		}
	}

	// Добавим недостающие элементы
	for i := 0; i < len(responsibilityIdList); i++ {
		if !Contains(generalResponsibilityIdList, responsibilityIdList[i]) {
			s.repo.AddWorkstationResponsibility(workstationId, responsibilityIdList[i])
		}
	}

	// Получим новые данные из БД о рабочих станциях
	listWorkstation, err := s.GetWorkstationList()
	if err != nil {
		logrus.Print(err)
		return nil, err
	}

	return listWorkstation, nil
}

// UpdateWorkstationStatus обновляет рабочее место
func (s *WorkstationService) UpdateWorkstationStatus(workstationId int, workstationStatus bool, employeeId int) (sql.Result, error) {
	return s.repo.UpdateWorkstationStatus(workstationId, workstationStatus, employeeId)
}
