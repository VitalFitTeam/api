package instructorhandler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/mailer"
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
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err = h.services.InstructorServices.CreateInstructor(ctx, instructor, hashToken)
	if err != nil {
		switch err {
		case shared_errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	//send email -> error -> rollback
	status, err := h.services.AuthServices.MailSenderStaff(ctx, instructor.User, plainToken, mailer.UserStaffActivate)
	if err != nil {
		h.services.Logger.Errorw("error sending activation url to email", "error", err)
		if err := h.services.AuthServices.DeleteResetToken(ctx, instructor.UserID); err != nil {
			h.services.Logger.Errorw("error deleting user activation token ", "error", err)
			return
		}
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(status, gin.H{"message": "instructor created"})

}

// @Summary		List all instructors
// @Description	Retrieves a list of all instructors in the system.
// @Tags			Instructors
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=[]InstructorResponse}	"List of instructors"
// @Failure		500	{object}	map[string]interface{}				"Internal Server Error"
// @Router			/instructor [get]
func (h *InstructorHandlers) GetInstructorsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	instructors, err := h.services.InstructorServices.GetInstructors(ctx)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	response := make([]*InstructorResponse, 0, len(instructors))
	for _, instructor := range instructors {
		ins := &InstructorResponse{
			InstructorID:      instructor.InstructorID,
			UserID:            instructor.UserID,
			FirstName:         instructor.User.FirstName,
			LastName:          instructor.User.LastName,
			Email:             instructor.User.Email,
			Phone:             instructor.User.Phone,
			IdentityDocument:  instructor.User.IdentityDocument,
			BirthDate:         instructor.User.BirthDate,
			Gender:            string(instructor.User.Gender),
			ProfilePictureURL: instructor.User.ProfilePictureURL,
			Speciality:        instructor.Speciality,
			Biography:         instructor.Biography,
		}
		response = append(response, ins)
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
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
// @Param			id	path		string							true	"Instructor UUID"
// @Success		200	{object}	object{data=InstructorResponse}	"Instructor details"
// @Failure		400	{object}	map[string]interface{}			"Bad Request: Invalid UUID format"
// @Failure		404	{object}	map[string]interface{}			"Not Found: Instructor not found"
// @Failure		500	{object}	map[string]interface{}			"Internal Server Error"
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

	response := &InstructorResponse{
		InstructorID:      instructor.InstructorID,
		FirstName:         instructor.User.FirstName,
		LastName:          instructor.User.LastName,
		Email:             instructor.User.Email,
		Phone:             instructor.User.Phone,
		IdentityDocument:  instructor.User.IdentityDocument,
		BirthDate:         instructor.User.BirthDate,
		Gender:            string(instructor.User.Gender),
		ProfilePictureURL: instructor.User.ProfilePictureURL,
		Speciality:        instructor.Speciality,
		Biography:         instructor.Biography,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})

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
	instructor, err := payload.toInstructor(instructorID)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

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

// @Summary		Assign instructors to a branch
// @Description	Assigns one or more instructors to a specific branch using their UUIDs.
// @Tags			Branch Instructors
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id			path		string								true	"Branch UUID"
// @Param			instructors	body		AssignInstructorsToBranchPayload	true	"Payload with instructor UUIDs to assign"
// @Success		204			{object}	nil									"Instructors assigned successfully"
// @Failure		400			{object}	map[string]interface{}				"Bad Request: Invalid UUID or payload"
// @Failure		404			{object}	map[string]interface{}				"Not Found: Branch or one of the instructors not found"
// @Failure		409			{object}	map[string]interface{}				"Conflict: Instructor already assigned to this branch"
// @Failure		500			{object}	map[string]interface{}				"Internal Server Error"
// @Router			/branches/{id}/instructor [post]
func (h *InstructorHandlers) AssignInstructorsToBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload AssignInstructorsToBranchPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	instructor, err := payload.toInstructor()

	err = h.services.InstructorServices.AssignInstructorsToBranch(ctx, branchID, instructor)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)

}

// @Summary		List instructors in a branch
// @Description	Retrieves a list of all instructors assigned to a specific branch.
// @Tags			Branch Instructors
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string									true	"Branch UUID"
// @Success		200	{object}	object{data=[]BranchInstructorResponse}	"List of branch instructors"
// @Failure		400	{object}	map[string]interface{}					"Bad Request: Invalid UUID format"
// @Failure		404	{object}	map[string]interface{}					"Not Found: Branch not found"
// @Failure		500	{object}	map[string]interface{}					"Internal Server Error"
// @Router			/branches/{id}/instructor [get]
func (h *InstructorHandlers) ListBranchInstructorsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	instructors, err := h.services.InstructorServices.ListBranchInstructors(ctx, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	response := make([]*BranchInstructorResponse, 0, len(instructors))
	for _, instructor := range instructors {
		ins := &BranchInstructorResponse{
			InstructorID:   instructor.InstructorID,
			InstructorName: instructor.User.FirstName + " " + instructor.User.LastName,
			Email:          instructor.User.Email,
			Phone:          instructor.User.Phone,
		}
		response = append(response, ins)
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// @Summary		Remove an instructor from a branch
// @Description	Removes a specific instructor from a specific branch.
// @Tags			Branch Instructors
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id				path		string					true	"Branch UUID"
// @Param			instructor_id	path		string					true	"Instructor UUID"
// @Success		204				{object}	nil						"Instructor removed successfully"
// @Failure		400				{object}	map[string]interface{}	"Bad Request: Invalid UUID format"
// @Failure		404				{object}	map[string]interface{}	"Not Found: Branch, instructor, or assignment not found"
// @Failure		500				{object}	map[string]interface{}	"Internal Server Error"
// @Router			/branches/{id}/instructor/{instructor_id} [delete]
func (h *InstructorHandlers) RemoveInstructorFromBranchHandler(c *gin.Context) {
	ctx := c.Request.Context()
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	instructorID, err := uuid.Parse(c.Param("instructor_id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	err = h.services.InstructorServices.RemoveInstructorFromBranch(ctx, branchID, instructorID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)

}
