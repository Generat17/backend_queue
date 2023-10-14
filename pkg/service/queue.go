package service

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/mail.v2"
	"os"
	"server/pkg/repository"
	"server/types"
	"strconv"
	"time"
)

type QueueService struct {
	repo        repository.Queue
	checkUpdate *[1]int
}

func NewQueueService(repo repository.Queue, checkUpdate *[1]int) *QueueService {

	return &QueueService{repo: repo, checkUpdate: checkUpdate}
}

// GetQueueList возвращает список клиентов(элементов) в очереди
func (s *QueueService) GetQueueList() ([]types.QueueItem, error) {
	// получаем очередь
	queue, err := s.repo.GetQueue()

	/// Удалить лишние элементы перед отправкой на фронт
	for i := 0; i < len(queue); i++ {
		if queue[i].ServiceTime != -1 {
			copy(queue[i:], queue[i+1:])            // удаляем элемент из очереди
			queue[len(queue)-1] = types.QueueItem{} // удаляем элемент из очереди
			queue = queue[:len(queue)-1]            // удаляем элемент из очереди
		}
	}
	if err != nil {
		return []types.QueueItem{}, err
	}
	return queue, nil
}

// GetQueueAdminList возвращает список клиентов(элементов) в очереди
func (s *QueueService) GetQueueAdminList() ([]types.QueueItem, error) {
	// получаем очередь
	queue, err := s.repo.GetQueue()
	if err != nil {
		return []types.QueueItem{}, err
	}
	return queue, nil
}

// GetQueueItemStatus возвращает список клиентов(элементов) в очереди
func (s *QueueService) GetQueueItemStatus(workstationId int) (types.QueueItem, error) {
	// получаем очередь
	queue, err := s.repo.GetQueueItemStatus(workstationId)
	if err != nil {
		logrus.Print("error GetQueueItemStatus")
		return types.QueueItem{}, err
	}

	if len(queue) == 0 {
		return types.QueueItem{Id: -1, Priority: 0, ServiceType: "", WorkstationNumber: -1, Status: 0, StartTime: 0, CallTime: 0, ServiceTime: 0, Quality: -1}, nil
	} else {
		return queue[0], nil
	}
}

// AddQueueItem добавляет нового клиента (элемент) в конец очереди
func (s *QueueService) AddQueueItem(service string) (int, error) {

	// Добавляем новый элемент в очередь
	res, err := s.repo.AddItemQueue(service)

	if err != nil {
		logrus.Print("error AddQueueItem")
		logrus.Print(err)
	}

	return res, nil
}

// ConfirmClient подтверждает, что клиент подошел к рабочему месту сотрудника
func (s *QueueService) ConfirmClient(numberQueue, employeeId int) (int, error) {
	s.repo.SetStatusEmployee(3, employeeId)
	s.repo.SetStatusQueueItem(3, numberQueue)
	s.checkUpdate[0] = s.checkUpdate[0] + 1

	_, err := s.repo.ConfirmClientTime(numberQueue)
	if err != nil {
		return -1, err
	}

	return 3, nil
}

// NotCome клиент не пришел
func (s *QueueService) NotCome(numberQueue, employeeId int) (int, error) {
	s.repo.SetStatusEmployee(1, employeeId)

	_, err := s.repo.NotCome(numberQueue)
	if err != nil {
		return -1, err
	}

	return 1, nil
}

// EndClient завершает обслуживание клиента
func (s *QueueService) EndClient(numberQueue, employeeId int) (int, error) {

	_, err := s.repo.EndClientTime(numberQueue)
	if err != nil {
		return -1, err
	}

	s.repo.SetStatusEmployee(1, employeeId)
	s.checkUpdate[0] = s.checkUpdate[0] + 1

	return 1, nil
}

// SetEmployeeStatus устанавливает статус для сотрудника
func (s *QueueService) SetEmployeeStatus(statusCode, employeeId int) (bool, error) {
	s.repo.SetStatusEmployee(statusCode, employeeId)

	return true, nil
}

// GetNewClient Выбирает доступного клиента из очереди и возращает информацию о нем
func (s *QueueService) GetNewClient(employeeId, workstationId int) (types.GetNewClientResponse, error) {
	// получаем список обязанностей сотрудника
	responsibilityEmployeeList, err := s.repo.GetResponsibilityByEmployeeId(employeeId)
	if err != nil {
		return types.GetNewClientResponse{NumberTicket: -1, ServiceTicket: "", EmployeeStatus: 1, NumberQueue: 0}, err
	}

	// получаем список обязанностей сотрудника
	responsibilityWorkstationList, err := s.repo.GetResponsibilityByWorkstationId(workstationId)
	if err != nil {
		return types.GetNewClientResponse{NumberTicket: -1, ServiceTicket: "", EmployeeStatus: 1, NumberQueue: 0}, err
	}

	// находим общие обязанности у сотрудника и рабочего места
	var generalResponsibility = []types.Responsibility{}

	for i := 0; i < len(responsibilityEmployeeList); i++ {
		for j := 0; j < len(responsibilityWorkstationList); j++ {
			if responsibilityEmployeeList[i] == responsibilityWorkstationList[j] {
				generalResponsibility = append(generalResponsibility, responsibilityEmployeeList[i])
			}
		}
	}

	queue, _ := s.repo.GetQueue()

	currentClientId := -1       // в эту переменную запишем номер элемента, который будет вызван
	currentClientPriority := -1 // в эту переменную запишем значение приоритета, текущего выбранного элемента

	// Пробегаемся по очереди и смотрим кого можно принять
	for i := 0; i < len(queue); i++ {
		for j := 0; j < len(generalResponsibility); j++ {
			if (queue[i].Status == 1) && (queue[i].ServiceType == generalResponsibility[j].ResponsibilityName) && (generalResponsibility[j].ResponsibilityPriority > currentClientPriority) {
				currentClientId = i
				currentClientPriority = generalResponsibility[j].ResponsibilityPriority
			}
		}
	}

	if currentClientId != -1 {
		s.repo.SetStatusEmployee(2, employeeId)
		s.repo.SetStatusQueueItem(2, queue[currentClientId].Id)
		s.repo.CallClientTime(queue[currentClientId].Id, workstationId)
		return types.GetNewClientResponse{NumberTicket: queue[currentClientId].Id, ServiceTicket: queue[currentClientId].ServiceType, EmployeeStatus: 2, NumberQueue: currentClientId}, nil
	}

	return types.GetNewClientResponse{NumberTicket: -1, ServiceTicket: "Нет доступного клиента", EmployeeStatus: 1, NumberQueue: 0}, nil
}

// ClearQueue очищает
func (s *QueueService) ClearQueue() (bool, error) {
	_, err := s.repo.ClearQueue()
	if err != nil {
		logrus.Print(err)
		return false, nil
	}

	return true, nil
}

// ClearLog очищает
func (s *QueueService) ClearLog() (bool, error) {
	_, err := s.repo.ClearLog()
	if err != nil {
		logrus.Print(err)
		return false, nil
	}

	return true, nil
}

// UpdateQuality устанавливает
func (s *QueueService) UpdateQuality(quality, client int) (bool, error) {
	s.repo.UpdateQuality(quality, client)

	return true, nil
}

// GetClientsLog возвращает список клиентов(элементов) в очереди
func (s *QueueService) GetClientsLog() ([]types.LogItem, error) {
	// получаем логи
	logs, err := s.repo.GetLogs()

	if err != nil {
		logrus.Print(err)
		return []types.LogItem{}, err
	}

	return logs, nil
}

// SendEmail возвращает список клиентов(элементов) в очереди
func (s *QueueService) SendEmail(to []string, subject, body string) (bool, error) {

	// Получение данных из конфигов и env файлова
	from := viper.GetString("email.from")
	host := viper.GetString("email.host")
	port, errPort := strconv.Atoi(viper.GetString("email.port"))
	if errPort != nil {
		logrus.Print(errPort)
	}
	password := os.Getenv("EMAIL_PASSWORD")

	for i := 0; i < len(to); i++ {
		abc := gomail.NewMessage()

		abc.SetHeader("From", from)
		abc.SetHeader("To", to[i])
		abc.SetHeader("Subject", subject)
		abc.SetBody("text/plain", body)

		a := gomail.NewDialer(host, port, from, password)
		if err := a.DialAndSend(abc); err != nil {
			logrus.Print(err)
		}
	}

	// Вызов функции, тестирование
	//s.SendEmail([]string{"alievtm@gmail.com", "lavello.michel@yandex.ru"}, "Заголовок тест", "Тело сообщения")

	return true, nil
}

// CheckLongWait проверяет долгое ожидание
func (s *QueueService) CheckLongWait() (bool, error) {
	email, err := s.repo.GetEmail()
	if err != nil {
		logrus.Print(err)
	}

	timing, err := s.repo.GetTiming()
	if err != nil {
		logrus.Print(err)
	}

	queue, err := s.repo.GetQueue()
	if err != nil {
		logrus.Print(err)
	}

	// Получение текущего времени
	time, err := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		logrus.Print(err)
	}

	for i := 0; i < len(queue); i++ {
		waitTime := time - queue[i].StartTime

		subject := "Длительное ожидание очереди. Более " + strconv.Itoa(waitTime/60) + " минут"
		body := "Клиент с номером " + strconv.Itoa(queue[i].Id) + " ожидает оказание услуги \"" + queue[i].ServiceType + "\" более " + strconv.Itoa(waitTime/60) + " минут"

		for k := 0; k < len(timing); k++ {
			if k == len(timing)-1 {
				if waitTime >= timing[k].Seconds {
					for j := 0; j < len(email); j++ {
						if email[j].Timing == timing[k].Id {
							if timing[k].Active {
								s.SendEmail([]string{email[j].Email}, subject, body)
							}
						}
					}
				}
			} else {
				if waitTime >= timing[k].Seconds && waitTime < timing[k+1].Seconds {
					for j := 0; j < len(email); j++ {
						if email[j].Timing == timing[k].Id {
							if timing[k].Active {
								s.SendEmail([]string{email[j].Email}, subject, body)
							}
						}
					}
				}
			}
		}
	}

	return true, nil
}

// GetTiming возвращает список клиентов(элементов) в очереди
func (s *QueueService) GetTiming() ([]types.TimingItem, error) {
	// получаем логи
	timing, err := s.repo.GetTiming()

	if err != nil {
		logrus.Print(err)
		return []types.TimingItem{}, err
	}

	return timing, nil
}

// GetEmail возвращает список клиентов(элементов) в очереди
func (s *QueueService) GetEmail() ([]types.EmailItem, error) {
	// получаем логи
	email, err := s.repo.GetEmail()

	if err != nil {
		logrus.Print(err)
		return []types.EmailItem{}, err
	}

	return email, nil
}

// UpdateTiming обновляет обязанность
func (s *QueueService) UpdateTiming(id int, seconds int, name string) (sql.Result, error) {
	return s.repo.UpdateTiming(id, seconds, name)
}

// RemoveTiming удаляет обязанность
func (s *QueueService) RemoveTiming(id int) (sql.Result, error) {
	return s.repo.RemoveTiming(id)
}

// AddTiming создает обязанность
func (s *QueueService) AddTiming(seconds int, name string) (sql.Result, error) {
	return s.repo.AddTiming(seconds, name)
}

// RemoveEmail удаляет обязанность
func (s *QueueService) RemoveEmail(id int) (sql.Result, error) {
	return s.repo.RemoveEmail(id)
}

// AddEmail создает обязанность
func (s *QueueService) AddEmail(timing int, email string) (sql.Result, error) {
	return s.repo.AddEmail(timing, email)
}

// ActiveTiming обновляет обязанность
func (s *QueueService) ActiveTiming(id int) (sql.Result, error) {
	return s.repo.ActiveTiming(id)
}

// RestartIdentity hello
func (s *QueueService) RestartIdentity() (sql.Result, error) {
	return s.repo.RestartIdentity()
}
