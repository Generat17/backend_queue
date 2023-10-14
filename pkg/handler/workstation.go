package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/types"
	"strconv"
)

type getWorkstationListsResponse struct {
	Data []types.Workstation `json:"data"`
}

// @Summary Get Workstation Data
// @Tags workstation
// @Description get workstation data
// @ID get-workstation
// @Accept  json
// @Produce  json
// @Success 200 {object} getWorkstationListsResponse "response"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/workstation  [post]
func (h *Handler) getWorkstationList(c *gin.Context) {
	items, err := h.services.Workstation.GetWorkstationList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getWorkstationListsResponse{Data: items})
}

type addWorkstationInput struct {
	WorkstationName string `json:"workstationName" binding:"required"`
}

// @Summary Add Workstation
// @Tags workstation
// @Description add workstation
// @ID add-workstation
// @Accept  json
// @Produce  json
// @Param input body WorkstationName true "credentials"
// @Success 200 {object} types.ResponseWorkstation
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/workstation/add [post]
func (h *Handler) addWorkstation(c *gin.Context) {

	var input addWorkstationInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.services.Workstation.AddWorkstation(input.WorkstationName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ResponseWorkstation{Response: response})
}

type updateWorkstationInput struct {
	WorkstationId   string `json:"workstationId" binding:"required"`
	WorkstationName string `json:"workstationName" binding:"required"`
}

// @Summary Update Workstation
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
func (h *Handler) updateWorkstation(c *gin.Context) {
	var input updateWorkstationInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	workstationId, _ := strconv.Atoi(input.WorkstationId)

	response, err := h.services.Workstation.UpdateWorkstation(workstationId, input.WorkstationName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ResponseWorkstation{Response: response})
}

type removeWorkstationInput struct {
	WorkstationId string `json:"workstationId" binding:"required"`
}

// @Summary Remove Workstation
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
func (h *Handler) removeWorkstation(c *gin.Context) {
	var input removeWorkstationInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	workstationId, _ := strconv.Atoi(input.WorkstationId)

	response, err := h.services.Workstation.RemoveWorkstation(workstationId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ResponseWorkstation{Response: response})
}

type updateWorkstationResponsibilityInput struct {
	WorkstationId        string `json:"workstationId" binding:"required"`
	ResponsibilityIdList []int  `json:"responsibilityIdList" binding:"required"`
}

// @Summary Update Workstation-Responsibility
// @Tags workstation-responsibility
// @Description remove workstation-responsibility
// @ID remove-workstation-responsibility
// @Accept  json
// @Produce  json
// @Param input body removeWorkstationResponsibilityInput true "credentials"
// @Success 200 {object} types.ResponseWorkstation
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/workstationResponsibility/remove [post]
func (h *Handler) updateWorkstationResponsibility(c *gin.Context) {
	var input updateWorkstationResponsibilityInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	workstationId, _ := strconv.Atoi(input.WorkstationId)
	responsibilityId := input.ResponsibilityIdList

	items, err := h.services.Workstation.UpdateWorkstationResponsibility(workstationId, responsibilityId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getWorkstationListsResponse{Data: items})
}
