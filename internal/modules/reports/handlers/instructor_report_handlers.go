package reporthandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Get Instructor Next Class
// @Description	Retrieves the start time of the next class for a specific instructor today.
// @Tags			Reports Instructor
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=string}		"Next class time or 'Sin pendientes'"
// @Failure		400	{object}	object{error=string}	"Bad Request"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/reports/instructors/next-class [get]
func (h *ReportHanlders) GetInstructorNextClassHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := h.services.UserServices.GetUserFromContext(c)

	// Buscar el InstructorID asociado al UserID del token
	instructorID, err := h.services.ReportServices.GetInstructorIDByUserID(ctx, user.UserID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	result, err := h.services.ReportServices.GetInstructorNextClass(ctx, *instructorID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
