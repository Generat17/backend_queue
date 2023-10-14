package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/types"
	"strconv"
)

// @Summary Get Queue List
// @Tags queue
// @Description get all queue list
// @ID get-queue-list
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.QueueItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue [get]
func (h *Handler) getQueueLists(c *gin.Context) {
	items, err := h.services.Queue.GetQueueList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Get Queue List
// @Tags queue
// @Description get all queue list
// @ID get-queue-list
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.QueueItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue [get]
func (h *Handler) getQueueAdminList(c *gin.Context) {
	items, err := h.services.Queue.GetQueueAdminList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Add New Queue Item
// @Tags queue
// @Description add new ticket (item queue) in the end of the queue
// @ID add-new-ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} types.QueueItemNumber
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue/service [get]
func (h *Handler) addQueueItem(c *gin.Context) {
	serviceType := c.Param("service")

	queueItemNumber, err := h.services.Queue.AddQueueItem(serviceType)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.QueueItemNumber{Ticket: queueItemNumber})
}

// @Summary get Queue Item Status
// @Tags queue
// @Description add new ticket (item queue) in the end of the queue
// @ID add-new-ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} types.QueueItemNumber
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue/service [get]
func (h *Handler) getQueueItemStatus(c *gin.Context) {
	workstationId, _ := strconv.Atoi(c.Param("workstation"))

	queueItem, err := h.services.Queue.GetQueueItemStatus(workstationId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, queueItem)
}

// @Summary update Quality
// @Tags queue
// @Description add new ticket (item queue) in the end of the queue
// @ID add-new-ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} types.QueueItemNumber
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue/service [get]
func (h *Handler) updateQuality(c *gin.Context) {
	client, _ := strconv.Atoi(c.Param("client"))
	quality, _ := strconv.Atoi(c.Param("quality"))

	res, err := h.services.Queue.UpdateQuality(quality, client)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Clear Queue
// @Tags queue
// @Description add new ticket (item queue) in the end of the queue
// @ID add-new-ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} types.QueueItemNumber
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue/service [get]
func (h *Handler) clearQueue(c *gin.Context) {

	res, err := h.services.Queue.ClearQueue()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Clear Queue
// @Tags queue
// @Description add new ticket (item queue) in the end of the queue
// @ID add-new-ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} types.QueueItemNumber
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue/service [get]
func (h *Handler) clearLog(c *gin.Context) {

	res, err := h.services.Queue.ClearLog()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary get Clients Log
// @Tags queue
// @Description add new ticket (item queue) in the end of the queue
// @ID add-new-ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} types.QueueItemNumber
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue/service [get]
func (h *Handler) getClientsLog(c *gin.Context) {

	res, err := h.services.Queue.GetClientsLog()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Get Queue List
// @Tags queue
// @Description get all queue list
// @ID get-queue-list
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.QueueItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue [get]
func (h *Handler) getTimingList(c *gin.Context) {
	items, err := h.services.Queue.GetTiming()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Get Queue List
// @Tags queue
// @Description get all queue list
// @ID get-queue-list
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.QueueItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue [get]
func (h *Handler) getEmailList(c *gin.Context) {
	items, err := h.services.Queue.GetEmail()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

type addTimingInput struct {
	Name    string `json:"name" binding:"required"`
	Seconds string `json:"seconds" binding:"required"`
}

// @Summary Add Responsibility
// @Tags responsibility
// @Description add responsibility
// @ID add-responsibility
// @Accept  json
// @Produce  json
// @Param input body addResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseResponsibility
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/responsibility/add [post]
func (h *Handler) addTiming(c *gin.Context) {

	var input addTimingInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	seconds, _ := strconv.Atoi(input.Seconds)

	response, err := h.services.Queue.AddTiming(seconds, input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

type updateTimingInput struct {
	Id      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Seconds string `json:"seconds" binding:"required"`
}

// @Summary Update Responsibility
// @Tags responsibility
// @Description update responsibility
// @ID update-responsibility
// @Accept  json
// @Produce  json
// @Param input body updateResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseResponsibility
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/responsibility/update [post]
func (h *Handler) updateTiming(c *gin.Context) {
	var input updateTimingInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := strconv.Atoi(input.Id)
	seconds, _ := strconv.Atoi(input.Seconds)

	response, err := h.services.Queue.UpdateTiming(id, seconds, input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

type removeTimingInput struct {
	Id string `json:"id" binding:"required"`
}

// @Summary Remove Responsibility
// @Tags responsibility
// @Description remove responsibility
// @ID remove-responsibility
// @Accept  json
// @Produce  json
// @Param input body removeResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseResponsibility
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/responsibility/remove [post]
func (h *Handler) removeTiming(c *gin.Context) {
	var input removeTimingInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := strconv.Atoi(input.Id)

	response, err := h.services.Queue.RemoveTiming(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

type addEmailInput struct {
	Timing string `json:"timing" binding:"required"`
	Email  string `json:"email" binding:"required"`
}

// @Summary Add Responsibility
// @Tags responsibility
// @Description add responsibility
// @ID add-responsibility
// @Accept  json
// @Produce  json
// @Param input body addResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseResponsibility
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/responsibility/add [post]
func (h *Handler) addEmail(c *gin.Context) {

	var input addEmailInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	timing, _ := strconv.Atoi(input.Timing)

	response, err := h.services.Queue.AddEmail(timing, input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

type removeEmailInput struct {
	Id string `json:"id" binding:"required"`
}

// @Summary Remove Responsibility
// @Tags responsibility
// @Description remove responsibility
// @ID remove-responsibility
// @Accept  json
// @Produce  json
// @Param input body removeResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseResponsibility
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/responsibility/remove [post]
func (h *Handler) removeEmail(c *gin.Context) {
	var input removeEmailInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := strconv.Atoi(input.Id)

	response, err := h.services.Queue.RemoveEmail(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

type activeTimingInput struct {
	Id string `json:"id" binding:"required"`
}

// @Summary Remove Responsibility
// @Tags responsibility
// @Description remove responsibility
// @ID remove-responsibility
// @Accept  json
// @Produce  json
// @Param input body removeResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseResponsibility
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/responsibility/remove [post]
func (h *Handler) activeTiming(c *gin.Context) {
	var input activeTimingInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := strconv.Atoi(input.Id)

	response, err := h.services.Queue.ActiveTiming(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Clear Queue
// @Tags queue
// @Description add new ticket (item queue) in the end of the queue
// @ID add-new-ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} types.QueueItemNumber
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/queue/service [get]
func (h *Handler) restartIdentity(c *gin.Context) {

	res, err := h.services.Queue.RestartIdentity()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
