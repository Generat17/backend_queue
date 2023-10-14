package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"server/pkg/repository"
	"server/types"
	"strconv"
	"time"
)

type tokenClaimsWorkstation struct {
	jwt.StandardClaims
	UserId        int `json:"user_id"`
	WorkstationId int `json:"workstation_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateEmployee создает сотрудника
func (s *AuthService) CreateEmployee(username, password, firstName, secondName string, isAdmin bool) (int, error) {
	passwordHash := generatePasswordHash(password)
	return s.repo.CreateEmployee(username, passwordHash, firstName, secondName, isAdmin)
}

// GenerateTokenWorkstation создает jwt access token
func (s *AuthService) GenerateTokenWorkstation(username, password string, workstationId int) (string, error) {
	employee, err := s.repo.GetEmployeeId(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	accessTokenTTLString, _ := strconv.Atoi(viper.GetString("token.accessTokenTTL"))
	accessTokenTTL := time.Duration(accessTokenTTLString) * time.Second

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaimsWorkstation{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		employee.EmployeeId,
		workstationId,
	})

	return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}

// GetEmployee получает данные сотрудника по логину и паролю сотрудника
func (s *AuthService) GetEmployee(username, password string) (types.Employee, error) {
	employee, err := s.repo.GetEmployee(username, generatePasswordHash(password))
	if err != nil {
		return types.Employee{EmployeeId: 0, Username: "", Password: "", FirstName: "", SecondName: "", SessionState: false, Status: 0}, err
	}

	return employee, nil
}

// GetEmployeeById получает данные сотрудника по ID сотрудника
func (s *AuthService) GetEmployeeById(employeeId int) (types.Employee, error) {
	employee, err := s.repo.GetEmployeeById(employeeId)
	if err != nil {
		return types.Employee{EmployeeId: 0, Username: "", Password: "", FirstName: "", SecondName: "", SessionState: false, Status: 0}, err
	}

	return employee, nil
}

// UpdateTokenWorkstation проверяет возможность выдачи нового jwt access token, при успехе, выдает новый jwt access token
func (s *AuthService) UpdateTokenWorkstation(employeeId, workstationId int, refreshToken string) (string, error) {

	var getSessionInfo types.SessionInfo
	getSessionInfo, err := s.repo.GetSession(employeeId)
	if err != nil {
		return "error check session", err
	}

	timeNow := time.Now().Unix()

	refreshTokenTTL, err := strconv.ParseInt(viper.GetString("token.refreshTokenTTL"), 10, 64)
	if err != nil {
		return "error convert refreshTokenTTL", err
	}

	accessTokenTTLString, _ := strconv.Atoi(viper.GetString("token.accessTokenTTL"))
	accessTokenTTL := time.Duration(accessTokenTTLString) * time.Second

	if (refreshToken == getSessionInfo.RefreshToken) && (timeNow-getSessionInfo.ExpiresAt < refreshTokenTTL) && (getSessionInfo.Workstation == workstationId) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaimsWorkstation{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			employeeId,
			workstationId,
		})

		return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	}

	return "refreshToken is invalid", errors.New("refreshToken is invalid")
}

// ParseTokenWorkstation парсит jwt access token
func (s *AuthService) ParseTokenWorkstation(accessToken string) (types.ParseTokenWorkstationResponse, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaimsWorkstation{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return types.ParseTokenWorkstationResponse{UserId: 0, WorkstationId: 0}, err
	}

	claims, ok := token.Claims.(*tokenClaimsWorkstation)
	if !ok {
		return types.ParseTokenWorkstationResponse{UserId: 0, WorkstationId: 0}, errors.New("token claims are not of type *tokenClaims")
	}

	return types.ParseTokenWorkstationResponse{UserId: claims.UserId, WorkstationId: claims.WorkstationId}, nil
}

// GenerateRefreshToken создает refresh token
func (s *AuthService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	str := rand.NewSource(time.Now().Unix())
	r := rand.New(str)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

// SetSession устанавливает состояние сессии сотрудника
func (s *AuthService) SetSession(refreshToken string, workstationId int, employeeId int) (bool, error) {

	timeNow := time.Now().Unix()

	_, err := s.repo.SetSession(refreshToken, timeNow, workstationId, employeeId)

	if err != nil {
		return false, errors.New("error set session")
	}

	return true, nil
}

// LogOut завершение сессии сотрудника
func (s *AuthService) LogOut(employeeId int) (bool, error) {

	logrus.Print(employeeId)
	_, err := s.repo.ClearSession(employeeId)

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetStatusEmployee Получение текущего статуса сотрудника
func (s *AuthService) GetStatusEmployee(employeeId int) (int, error) {

	statusEmployee, err := s.repo.GetStatusEmployee(employeeId)
	if err != nil {
		return -1, errors.New("error get status employee")
	}

	return statusEmployee, nil
}

// generatePasswordHash Хэширует пароль
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}
