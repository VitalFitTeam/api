package routinehandlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	routinedomain "github.com/vitalfit/api/internal/modules/routines/domain"
	"github.com/vitalfit/api/pkg/pagination"
)

// @Summary		Create a new routine template
// @Description	Creates a routine template with exercises. Only for Instructors/Admins.
// @Tags			Routines
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateRoutineRequest	true	"Routine creation payload"
// @Success		201		{object}	routinedomain.Routine	"Routine created successfully"
// @Failure		400		{object}	object{error=string}	"Bad Request"
// @Failure		500		{object}	object{error=string}	"Internal Server Error"
// @Router			/routines [post]
func (h *RoutineHandlers) CreateRoutineHandler(c *gin.Context) {
	var payload CreateRoutineRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user := h.services.UserServices.GetUserFromContext(c)

	// Map DTO to Domain
	routine := &routinedomain.Routine{
		Name:        payload.Name,
		Description: payload.Description,
		Level:       routinedomain.RoutineLevel(payload.Level),
		ServiceID:   payload.ServiceID,
		CreatorID:   &user.UserID,
	}

	for _, exDto := range payload.Exercises {
		routine.RoutineExercises = append(routine.RoutineExercises, routinedomain.RoutineExercise{
			ExerciseID: exDto.ExerciseID,
			Sets:       exDto.Sets,
			Reps:       exDto.Reps,
			RestTime:   exDto.RestTime,
			Order:      exDto.Order,
			Notes:      exDto.Notes,
		})
	}

	if err := h.services.Routine.CreateRoutine(c.Request.Context(), routine); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, routine)
}

// @Summary		Assign a routine to a client
// @Description	Assigns an existing routine template to a client. Only for Instructors.
// @Tags			Routines
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			payload	body		AssignRoutineRequest	true	"Assignment payload"
// @Success		201		{object}	map[string]string		"Success message"
// @Failure		400		{object}	object{error=string}	"Bad Request"
// @Failure		404		{object}	object{error=string}	"Not Found (Instructor/Client/Routine)"
// @Failure		500		{object}	object{error=string}	"Internal Server Error"
// @Router			/routines/assign [post]
func (h *RoutineHandlers) AssignRoutineHandler(c *gin.Context) {
	var payload AssignRoutineRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user := h.services.UserServices.GetUserFromContext(c)

	err := h.services.Routine.AssignRoutine(c.Request.Context(), user.UserID, payload.ClientID, payload.RoutineID, payload.DueDate)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Routine assigned successfully"})
}

// @Summary		Get my assigned routines
// @Description	Retrieves the list of active routines assigned to the authenticated client.
// @Tags			Routines
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	[]routinedomain.UserRoutine	"List of assigned routines"
// @Failure		500	{object}	object{error=string}		"Internal Server Error"
// @Router			/routines/my-routines [get]
func (h *RoutineHandlers) GetMyRoutinesHandler(c *gin.Context) {
	user := h.services.UserServices.GetUserFromContext(c)

	routines, err := h.services.Routine.GetClientRoutines(c.Request.Context(), user.UserID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, routines)
}

// @Summary		Create a new exercise
// @Description	Creates a new exercise in the library. Only for Instructors/Admins.
// @Tags			Routines
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			payload	body		CreateExerciseRequest	true	"Exercise creation payload"
// @Success		201		{object}	routinedomain.Exercise	"Exercise created successfully"
// @Failure		400		{object}	object{error=string}	"Bad Request"
// @Failure		500		{object}	object{error=string}	"Internal Server Error"
// @Router			/exercises [post]
func (h *RoutineHandlers) CreateExerciseHandler(c *gin.Context) {
	var payload CreateExerciseRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	exercise := &routinedomain.Exercise{
		Name:        payload.Name,
		Description: payload.Description,
		VideoURL:    payload.VideoURL,
		MuscleGroup: routinedomain.MuscleGroup(payload.MuscleGroup),
	}

	if err := h.services.Routine.CreateExercise(c.Request.Context(), exercise); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

// @Summary		List exercises
// @Description	Retrieves a paginated list of exercises from the catalog.
// @Tags			Routines
// @Security		ApiKeyAuth
// @Produce		json
// @Param			limit	query		int										false	"Number of results per page"
// @Param			page	query		int										false	"Page number"
// @Param			sort	query		string									false	"Sort order (asc/desc)"
// @Param			search	query		string									false	"Search term"
// @Success		200		{object}	object{data=[]routinedomain.Exercise}	"List of exercises"
// @Failure		500		{object}	object{error=string}					"Internal Server Error"
// @Router			/exercises [get]
func (h *RoutineHandlers) GetExercisesHandler(c *gin.Context) {
	fq := pagination.PaginatedFeedQuery{
		Limit:  10,
		Page:   1,
		Sort:   "asc",
		Search: "",
	}
	fq, err := fq.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	exercises, total, err := h.services.Routine.GetExercises(c.Request.Context(), fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	nextURL := fmt.Sprintf("/exercises?limit=%d&page=%d&sort=%s&search=%s", fq.Limit, fq.Page+1, fq.Sort, fq.Search)
	previousPage := fq.Page - 1
	if previousPage < 1 {
		previousPage = 1
	}
	previousURL := fmt.Sprintf("/exercises?limit=%d&page=%d&sort=%s&search=%s", fq.Limit, previousPage, fq.Sort, fq.Search)

	c.JSON(http.StatusOK, pagination.PaginatedResponseTotal[*routinedomain.Exercise]{
		Data:     exercises,
		Count:    int64(len(exercises)),
		Total:    total,
		Next:     nextURL,
		Previous: previousURL,
	})
}

// @Summary		Get client routines
// @Description	Retrieves the list of active routines assigned to a specific client. Only for Instructors/Admins.
// @Tags			Routines
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string						true	"Client ID (UUID)"
// @Success		200	{object}	[]routinedomain.UserRoutine	"List of assigned routines"
// @Failure		400	{object}	object{error=string}		"Bad Request"
// @Failure		500	{object}	object{error=string}		"Internal Server Error"
// @Router			/routines/client/{id} [get]
func (h *RoutineHandlers) GetClientRoutinesHandler(c *gin.Context) {
	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	routines, err := h.services.Routine.GetClientRoutines(c.Request.Context(), clientID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, routines)
}
