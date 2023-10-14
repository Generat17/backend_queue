package handler

import (
	"net/http"
	"server/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

type signUpInput struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	FirstName  string `json:"firstName" binding:"required"`
	SecondName string `json:"secondName" binding:"required"`
	IsAdmin    string `json:"isAdmin" binding:"required"`
}

// @Summary Sign Up
// @Tags auth
// @Description registration new account
// @ID registration-account
// @Accept  json
// @Produce  json
// @Param input body types.Employee true "credentials"
// @Success 200 {object} types.SignUpResponse "Employee ID"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var isAdmin bool
	if input.IsAdmin == "false" {
		isAdmin = false
	} else {
		isAdmin = true
	}

	id, err := h.services.Authorization.CreateEmployee(input.Username, input.Password, input.FirstName, input.SecondName, isAdmin)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.SignUpResponse{Id: id})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Sign In
// @Tags auth
// @Description employee authorization
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {object} types.AuthorizationResponse "response"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signInWorkstation(c *gin.Context) {
	var input signInInput
	workstationId, _ := strconv.Atoi(c.Param("workstation"))

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	workstation, err := h.services.Workstation.GetWorkstation(workstationId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	employee, err := h.services.Authorization.GetEmployee(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if workstation.IsBusy == true && workstation.EmployeeId != employee.EmployeeId {
		c.JSON(http.StatusOK, "workstation is busy")
		return
	}

	accessToken, err := h.services.Authorization.GenerateTokenWorkstation(input.Username, input.Password, workstationId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := h.services.Authorization.GenerateRefreshToken()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.services.Authorization.SetSession(refreshToken, workstationId, employee.EmployeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.services.Queue.SetEmployeeStatus(1, employee.EmployeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.services.Workstation.UpdateWorkstationStatus(workstation.WorkstationId, true, employee.EmployeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.AuthorizationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Employee:     employee,
		Workstation:  workstation,
	})
}

type refreshTokenInput struct {
	WorkstationId string `json:"workstationId" binding:"required"`
	EmployeeId    string `json:"employeeId" binding:"required"`
	RefreshToken  string `json:"refreshToken" binding:"required"`
}

// @Summary Refresh
// @Security ApiKeyAuth
// @Tags auth
// @Description refresh AccessToken, RefreshToken, Employee, Workstation data
// @ID refresh
// @Accept  json
// @Produce  json
// @Param input body refreshTokenInput true "credentials"
// @Success 200 {object} types.AuthorizationResponse "response"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {

	var input refreshTokenInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	WorkstationId, _ := strconv.Atoi(input.WorkstationId)
	EmployeeId, _ := strconv.Atoi(input.EmployeeId)

	accessToken, err := h.services.Authorization.UpdateTokenWorkstation(EmployeeId, WorkstationId, input.RefreshToken)
	if err != nil {
		if accessToken == "refreshToken is invalid" {
			newErrorResponse(c, http.StatusForbidden, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	workstation, err := h.services.Workstation.GetWorkstation(WorkstationId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := h.services.Authorization.GenerateRefreshToken()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	employee, err := h.services.Authorization.GetEmployeeById(EmployeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.services.Authorization.SetSession(refreshToken, WorkstationId, EmployeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.AuthorizationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Employee:     employee,
		Workstation:  workstation,
	})
}

type logoutInput struct {
	EmployeeId    string `json:"employeeId" binding:"required"`
	WorkstationId string `json:"workstationId" binding:"required"`
}

// @Summary Logout
// @Security ApiKeyAuth
// @Tags auth
// @Description logout account
// @ID logout
// @Accept  json
// @Produce  json
// @Param input body logoutInput true "credentials"
// @Success 200 {object} types.LogoutResponse "response"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	var input logoutInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	employeeId, _ := strconv.Atoi(input.EmployeeId)
	workstationId, _ := strconv.Atoi(input.WorkstationId)

	statusResponse, err := h.services.Authorization.LogOut(employeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.services.Workstation.UpdateWorkstationStatus(workstationId, false, -1)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.LogoutResponse{StatusResponse: statusResponse})
}
