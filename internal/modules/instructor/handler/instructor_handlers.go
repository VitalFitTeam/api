package instructorhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
)

// @Summary		Create a new instructor
// @Description	Creates a new instructor, which also creates an associated user with the 'instructor' role.
// @Tags			Instructors
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			instructor	body		CreateInstructorPayload	true	"Instructor creation payload"
// @Success		201			{object}	nil						"Instructor created successfully"
// @Failure		400			{object}	map[string]interface{}	"Bad Request: Invalid payload"
// @Failure		409			{object}	map[string]interface{}	"Conflict: User with this email or identity document already exists"
// @Failure		500			{object}	map[string]interface{}	"Internal Server Error"
// @Router			/instructor [post]
func (h *InstructorHandlers) CreateInstructorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload CreateInstructorPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	instructor, err := payload.toInstructor()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	err = h.services.InstructorServices.CreateInstructor(ctx, instructor)
	if err != nil {
		switch err {
		case shared_errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "instructor created"})

}

// @Summary		List all instructors
// @Description	Retrieves a list of all instructors in the system.
// @Tags			Instructors
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=[]instructordomain.Instructor}	"List of instructors"
// @Failure		500	{object}	map[string]interface{}						"Internal Server Error"
// @Router			/instructor [get]
func (h *InstructorHandlers) GetInstructorsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	instructors, err := h.services.InstructorServices.GetInstructors(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": instructors})
}

// @Summary		Delete an instructor
// @Description	Deletes a specific instructor by their UUID.
// @Tags			Instructors
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string					true	"Instructor UUID"
// @Success		204	{object}	nil						"Instructor deleted successfully"
// @Failure		400	{object}	map[string]interface{}	"Bad Request: Invalid UUID format"
// @Failure		404	{object}	map[string]interface{}	"Not Found: Instructor not found"
// @Failure		500	{object}	map[string]interface{}	"Internal Server Error"
// @Router			/instructor/{id} [delete]
func (h *InstructorHandlers) DeleteInstructorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	instructorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.InstructorServices.DeleteInstructor(ctx, instructorID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Get instructor by ID
// @Description	Retrieves detailed information about a specific instructor by their UUID.
// @Tags			Instructors
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string										true	"Instructor UUID"
// @Success		200	{object}	object{data=instructordomain.Instructor}	"Instructor details"
// @Failure		400	{object}	map[string]interface{}						"Bad Request: Invalid UUID format"
// @Failure		404	{object}	map[string]interface{}						"Not Found: Instructor not found"
// @Failure		500	{object}	map[string]interface{}						"Internal Server Error"
// @Router			/instructor/{id} [get]
func (h *InstructorHandlers) GetInstructorByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	instructorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	instructor, err := h.services.InstructorServices.GetInstructorByID(ctx, instructorID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": instructor})

}

// @Summary		Update an instructor
// @Description	Updates an instructor's speciality and biography.
// @Tags			Instructors
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id			path		string					true	"Instructor UUID"
// @Param			instructor	body		UpdateInstructorPayload	true	"Payload with fields to update"
// @Success		204			{object}	nil						"Instructor updated successfully"
// @Failure		400			{object}	map[string]interface{}	"Bad Request: Invalid UUID or payload"
// @Failure		404			{object}	map[string]interface{}	"Not Found: Instructor not found"
// @Failure		500			{object}	map[string]interface{}	"Internal Server Error"
// @Router			/instructor/{id} [put]
func (h *InstructorHandlers) UpdateInstructorHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload UpdateInstructorPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	instructorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	instructor := payload.toInstructor(instructorID)

	err = h.services.InstructorServices.UpdateInstructor(ctx, instructor)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)

}
