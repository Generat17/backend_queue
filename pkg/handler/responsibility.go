package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/types"
	"strconv"
)

// @Summary Get Responsibility List
// @Tags responsibility
// @Description get responsibility list
// @ID get-responsibility-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} []types.Responsibility
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/responsibility [get]
func (h *Handler) getResponsibilityList(c *gin.Context) {
	items, err := h.services.Responsibility.GetResponsibilityList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)
}

type addResponsibilityInput struct {
	ResponsibilityName string `json:"responsibilityName" binding:"required"`
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
func (h *Handler) addResponsibility(c *gin.Context) {

	var input addResponsibilityInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.services.Responsibility.AddResponsibility(input.ResponsibilityName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ResponseResponsibility{Response: response})
}

type updateResponsibilityInput struct {
	ResponsibilityId       string `json:"responsibilityId" binding:"required"`
	ResponsibilityName     string `json:"responsibilityName" binding:"required"`
	ResponsibilityPriority string `json:"ResponsibilityPriority" binding:"required"`
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
func (h *Handler) updateResponsibility(c *gin.Context) {
	var input updateResponsibilityInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	responsibilityId, _ := strconv.Atoi(input.ResponsibilityId)
	responsibilityPriority, _ := strconv.Atoi(input.ResponsibilityPriority)

	response, err := h.services.Responsibility.UpdateResponsibility(responsibilityId, input.ResponsibilityName, responsibilityPriority)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ResponseResponsibility{Response: response})
}

type removeResponsibilityInput struct {
	ResponsibilityId string `json:"responsibilityId" binding:"required"`
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
func (h *Handler) removeResponsibility(c *gin.Context) {
	var input removeResponsibilityInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	responsibilityId, _ := strconv.Atoi(input.ResponsibilityId)

	response, err := h.services.Responsibility.RemoveResponsibility(responsibilityId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, types.ResponseResponsibility{Response: response})
}
