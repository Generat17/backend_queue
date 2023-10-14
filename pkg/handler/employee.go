package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/types"
	"strconv"
)

// @Summary Get Employee Lists
// @Tags employee
// @Description get list employee
// @ID get-employee-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} types.GetEmployeeListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/employee [get]
func (h *Handler) getEmployeeList(c *gin.Context) {
	items, err := h.services.Employee.GetEmployeeList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.GetEmployeeListsResponse{Data: items})
}

// @Summary Get Employee Status by workstationID
// @Tags employee
// @Description get list employee
// @ID get-employee-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} types.GetEmployeeListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/employee [get]
func (h *Handler) getEmployeeStatus(c *gin.Context) {

	workstationId, _ := strconv.Atoi(c.Param("workstation"))

	employee, err := h.services.Employee.GetEmployeeStatus(workstationId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, employee)
}

// @Summary Get New Client
// @Security ApiKeyAuth
// @Tags client
// @Description get an available client from the queue
// @ID get-new-client
// @Accept  json
// @Produce  json
// @Success 200 {object} types.GetNewClientResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/employee/client [post]
func (h *Handler) getNewClient(c *gin.Context) {
	employeeId, _ := c.Get(userCtx)
	empId := employeeId.(int)
	workstationId, _ := c.Get(workstationCtx)
	workId := workstationId.(int)

	client, err := h.services.Queue.GetNewClient(empId, workId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, client)
}

type confirmClientInput struct {
	NumberQueue string `json:"numberQueue" binding:"required"`
}

// @Summary Confirm Client
// @Security ApiKeyAuth
// @Tags client
// @Description confirms that the client has approached the workstation
// @ID confirm-client
// @Accept  json
// @Produce  json
// @Param input body confirmClientInput true "credentials"
// @Success 200 {object} types.ConfirmClientResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/employee/confirmClient [post]
func (h *Handler) confirmClient(c *gin.Context) {
	employeeId, _ := c.Get(userCtx)
	empId := employeeId.(int)

	var input confirmClientInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	numberQueue, _ := strconv.Atoi(input.NumberQueue)

	client, err := h.services.Queue.ConfirmClient(numberQueue, empId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ConfirmClientResponse{NumberQueue: client})
}

type notComeClientInput struct {
	NumberQueue string `json:"numberQueue" binding:"required"`
}

// @Summary End Client
// @Security ApiKeyAuth
// @Tags client
// @Description complete the client
// @ID end-client
// @Accept  json
// @Produce  json
// @Success 200 {object} types.ConfirmClientResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/employee/endClient [post]
func (h *Handler) notComeClient(c *gin.Context) {
	employeeId, _ := c.Get(userCtx)
	empId := employeeId.(int)

	var input notComeClientInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	numberQueue, _ := strconv.Atoi(input.NumberQueue)

	client, err := h.services.Queue.NotCome(numberQueue, empId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ConfirmClientResponse{NumberQueue: client})
}

type endClientInput struct {
	NumberQueue string `json:"numberQueue" binding:"required"`
}

// @Summary End Client
// @Security ApiKeyAuth
// @Tags client
// @Description complete the client
// @ID end-client
// @Accept  json
// @Produce  json
// @Success 200 {object} types.ConfirmClientResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/employee/endClient [post]
func (h *Handler) endClient(c *gin.Context) {
	employeeId, _ := c.Get(userCtx)
	empId := employeeId.(int)

	var input endClientInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	numberQueue, _ := strconv.Atoi(input.NumberQueue)

	client, err := h.services.Queue.EndClient(numberQueue, empId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ConfirmClientResponse{NumberQueue: client})
}

// @Summary Get Employee Status
// @Security ApiKeyAuth
// @Tags employee
// @Description get the current status of an employee
// @ID get-status-employee
// @Accept  json
// @Produce  json
// @Success 200 {object} types.EmployeeStatusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/employee/getStatus [post]
func (h *Handler) getStatus(c *gin.Context) {
	employeeId, _ := c.Get(userCtx)
	empId := employeeId.(int)

	status, err := h.services.Authorization.GetStatusEmployee(empId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.EmployeeStatusResponse{EmployeeStatus: status})
}

type getEmployeeListsResponse struct {
	Data []types.Employee `json:"data"`
}

type updateEmployeeResponsibilityInput struct {
	EmployeeId           string `json:"employeeId" binding:"required"`
	ResponsibilityIdList []int  `json:"responsibilityIdList" binding:"required"`
}

// @Summary Update Employee-Responsibility
// @Tags employee-responsibility
// @Description remove employee-responsibility
// @ID remove-employee-responsibility
// @Accept  json
// @Produce  json
// @Param input body removeWorkstationResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseWorkstation
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/workstationResponsibility/remove [post]
func (h *Handler) updateEmployeeResponsibility(c *gin.Context) {
	var input updateEmployeeResponsibilityInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	employeeId, _ := strconv.Atoi(input.EmployeeId)
	responsibilityId := input.ResponsibilityIdList

	items, err := h.services.Employee.UpdateEmployeeResponsibility(employeeId, responsibilityId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getEmployeeListsResponse{Data: items})
}

type updateEmployeeInput struct {
	EmployeeId string `json:"employeeId" binding:"required"`
	Username   string `json:"username" binding:"required"`
	FirstName  string `json:"firstName" binding:"required"`
	SecondName string `json:"secondName" binding:"required"`
	IsAdmin    string `json:"isAdmin" binding:"required"`
}

// @Summary Update Employee
// @Tags workstation
// @Description update workstation
// @ID update-workstation
// @Accept  json
// @Produce  json
// @Param input body updateWorkstationInput true "credentials"
// @Success 200 {object} types.ResponseWorkstation
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/workstation/update [post]
func (h *Handler) updateEmployee(c *gin.Context) {
	var input updateEmployeeInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	employeeId, _ := strconv.Atoi(input.EmployeeId)

	var isAdmin bool
	if input.IsAdmin == "false" {
		isAdmin = false
	} else {
		isAdmin = true
	}

	response, err := h.services.Employee.UpdateEmployee(employeeId, input.Username, input.FirstName, input.SecondName, isAdmin)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

type removeEmployeeInput struct {
	EmployeeId string `json:"employeeId" binding:"required"`
}

// @Summary Remove Employee
// @Tags workstation
// @Description remove workstation
// @ID remove-workstation
// @Accept  json
// @Produce  json
// @Param input body removeWorkstationInput true "credentials"
// @Success 200 {object} types.ResponseWorkstation
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/workstation/remove [post]
func (h *Handler) removeEmployee(c *gin.Context) {
	var input removeEmployeeInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	employeeId, _ := strconv.Atoi(input.EmployeeId)

	response, err := h.services.Employee.RemoveEmployee(employeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
