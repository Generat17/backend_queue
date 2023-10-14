package service

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"server/pkg/repository"
	"server/types"
)

type EmployeeService struct {
	repo repository.Employee
}

func NewEmployeeService(repo repository.Employee) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// GetEmployeeList возвращает список сотрудников
func (s *EmployeeService) GetEmployeeList() ([]types.Employee, error) {
	listEmployeeResponsibility, err := s.repo.GetEmployeeResponsibilityList()
	if err != nil {
		return nil, err
	}

	listEmployee, err := s.repo.GetEmployeeList()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(listEmployee); i++ {

		for j := 0; j < len(listEmployeeResponsibility); j++ {
			if listEmployeeResponsibility[j].EmployeeId == listEmployee[i].EmployeeId {
				listEmployee[i].ResponsibilityList = append(listEmployee[i].ResponsibilityList, types.EmployeeResponsibilityItem{Id: listEmployeeResponsibility[j].ResponsibilityId, Name: listEmployeeResponsibility[j].ResponsibilityName})
			}
		}
	}

	return listEmployee, nil
}

// GetEmployeeStatus возвращает сотрудника и его статус
func (s *EmployeeService) GetEmployeeStatus(workstationId int) (types.EmployeeStatus, error) {
	employeeStatus, err := s.repo.GetEmployeeStatus(workstationId)
	if err != nil {
		logrus.Print("error get employee status")
	}

	if len(employeeStatus) == 0 {
		return types.EmployeeStatus{EmployeeId: -1, Status: -1}, nil
	} else {
		return employeeStatus[0], nil
	}
}

// GetEmployeeStatusList возвращает список сотрудников и их статус
func (s *EmployeeService) GetEmployeeStatusList() ([]types.EmployeeStatus, error) {
	employeeStatusList, err := s.repo.GetEmployeeStatusList()
	if err != nil {
		logrus.Print("error get employee status")
	}

	return employeeStatusList, nil
}

// UpdateEmployeeResponsibility обновляет связь сотрудник - обязанность
func (s *EmployeeService) UpdateEmployeeResponsibility(employeeId int, responsibilityIdList []int) ([]types.Employee, error) {
	currentEmployeeResponsibilityIdList, _ := s.repo.GetEmployeeResponsibilityById(employeeId) // получаем текущие связи

	var generalResponsibilityIdList = []int{} // в этом массиве храним общие функции рабочих мест

	// находим общие функции рабочих мест и записываем их в массив
	for i := 0; i < len(currentEmployeeResponsibilityIdList); i++ {
		for j := 0; j < len(responsibilityIdList); j++ {
			if currentEmployeeResponsibilityIdList[i] == responsibilityIdList[j] {
				generalResponsibilityIdList = append(generalResponsibilityIdList, currentEmployeeResponsibilityIdList[i])
			}
		}
	}

	// Удалим лишние элементы
	for i := 0; i < len(currentEmployeeResponsibilityIdList); i++ {
		if !Contains(generalResponsibilityIdList, currentEmployeeResponsibilityIdList[i]) {
			s.repo.RemoveEmployeeResponsibility(employeeId, currentEmployeeResponsibilityIdList[i])
		}
	}

	// Добавим недостающие элементы
	for i := 0; i < len(responsibilityIdList); i++ {
		if !Contains(generalResponsibilityIdList, responsibilityIdList[i]) {
			s.repo.AddEmployeeResponsibility(employeeId, responsibilityIdList[i])
		}
	}

	// Получим новые данные из БД о рабочих станциях
	listEmployee, err := s.GetEmployeeList()
	if err != nil {
		logrus.Print(err)
		return nil, err
	}

	return listEmployee, nil
}

// UpdateEmployee обновляет сотрудника
func (s *EmployeeService) UpdateEmployee(employeeId int, username, firstName, secondName string, isAdmin bool) (sql.Result, error) {
	return s.repo.UpdateEmployee(employeeId, username, firstName, secondName, isAdmin)
}

// RemoveEmployee удаляет сотрудника
func (s *EmployeeService) RemoveEmployee(employeeId int) (sql.Result, error) {
	return s.repo.RemoveEmployee(employeeId)
}
