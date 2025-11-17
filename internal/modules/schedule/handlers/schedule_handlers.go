package schedulehandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ------------------------------
// POST /branches/:id/schedule
// ------------------------------

// @Summary		Create a scheduled class
// @Description	Creates a class in the branch schedule
// @Tags			Schedule
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id		path		string					true	"Branch UUID"
// @Param			class	body		CreateClassPayload		true	"Class payload"
// @Success		201		{object}	map[string]interface{}	"Class created"
// @Failure		400		{object}	map[string]interface{}	"Invalid input"
// @Failure		500		{object}	map[string]interface{}	"Server error"
// @Router			/branches/{id}/schedule [post]
func (h *ScheduleHandlers) CreateClassHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var payload CreateClassPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	class, err := payload.ToClass(branchID)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.ScheduleServices.CreateClass(ctx, class); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Class created successfully",
		"class_id": class.ClassID,
	})
}

// ------------------------------
// GET /branches/:id/schedule
// ------------------------------

// @Summary		List classes for a branch
// @Description	Returns the scheduled classes for the branch
// @Tags			Schedule
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string	true	"Branch UUID"
// @Success		200	{object}	object{data=[]ClassResponse}
// @Failure		400	{object}	map[string]interface{}
// @Failure		500	{object}	map[string]interface{}
// @Router			/branches/{id}/schedule [get]
func (h *ScheduleHandlers) GetClassesByBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()

	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	classes, err := h.services.ScheduleServices.GetClassesByBranch(ctx, branchID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.services.LogErrors.NotFoundResponse(c)
			return
		default:
			h.services.LogErrors.InternalServerError(c, err)
			return
		}
	}

	resp := make([]ClassResponse, 0, len(classes))
	for _, class := range classes {
		resp = append(resp, ClassResponse{
			ClassID:      class.ClassID,
			BranchID:     class.BranchID,
			ServiceID:    class.ServiceID,
			InstructorID: class.InstructorID,
			StartsAt:     class.StartsAt,
			EndsAt:       class.EndsAt,
			MaxCapacity:  class.MaxCapacity,
			IsVisible:    class.IsVisible,
			Notes:        class.Notes,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// ------------------------------
// GET /schedule/:classId
// ------------------------------

// @Summary		Get class by ID
// @Description	Returns details of a scheduled class
// @Tags			Schedule
// @Security		ApiKeyAuth
// @Produce		json
// @Param			classId	path		string	true	"Class UUID"
// @Success		200		{object}	object{data=ClassResponse}
// @Failure		400		{object}	map[string]interface{}
// @Failure		500		{object}	map[string]interface{}
// @Router			/schedule/{classId} [get]
func (h *ScheduleHandlers) GetClassByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	classID, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	class, err := h.services.ScheduleServices.GetClassByID(ctx, classID)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.services.LogErrors.NotFoundResponse(c)
			return
		default:
			h.services.LogErrors.InternalServerError(c, err)
			return
		}
	}

	resp := ClassResponse{
		ClassID:      class.ClassID,
		BranchID:     class.BranchID,
		ServiceID:    class.ServiceID,
		InstructorID: class.InstructorID,
		StartsAt:     class.StartsAt,
		EndsAt:       class.EndsAt,
		MaxCapacity:  class.MaxCapacity,
		IsVisible:    class.IsVisible,
		Notes:        class.Notes,
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// ------------------------------
// PUT /schedule/:classId
// ------------------------------

// @Summary		Update scheduled class
// @Description	Updates a class in the calendar
// @Tags			Schedule
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			classId	path		string				true	"Class UUID"
// @Param			class	body		UpdateClassPayload	true	"Class update payload"
// @Success		204		{object}	nil
// @Failure		400		{object}	map[string]interface{}
// @Failure		500		{object}	map[string]interface{}
// @Router			/schedule/{classId} [put]
func (h *ScheduleHandlers) UpdateClassHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var payload UpdateClassPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	classID, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	class, err := payload.ToClassUpdate(classID)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.ScheduleServices.UpdateClass(ctx, class); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ------------------------------
// DELETE /schedule/:classId
// ------------------------------

// @Summary		Delete scheduled class
// @Description	Deletes or cancels a class from the calendar
// @Tags			Schedule
// @Security		ApiKeyAuth
// @Produce		json
// @Param			classId	path		string	true	"Class UUID"
// @Success		204		{object}	nil
// @Failure		400		{object}	map[string]interface{}
// @Failure		500		{object}	map[string]interface{}
// @Router			/schedule/{classId} [delete]
func (h *ScheduleHandlers) DeleteClassHandler(c *gin.Context) {
	ctx := c.Request.Context()

	classID, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.ScheduleServices.DeleteClass(ctx, classID); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
