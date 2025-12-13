package instructorhandler

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type InstructorHandlersInterface interface {
	InstructorRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	CreateInstructorHandler(c *gin.Context)
	GetInstructorsHandler(c *gin.Context)
	GetSummaryHandler(c *gin.Context)
	DeleteInstructorHandler(c *gin.Context)
	GetInstructorByIDHandler(c *gin.Context)
	UpdateInstructorHandler(c *gin.Context)

	AssignInstructorsToBranchHandler(c *gin.Context)
	ListBranchInstructorsHandler(c *gin.Context)
	RemoveInstructorFromBranchHandler(c *gin.Context)

	AssignInstructorSpecialtyHandler(c *gin.Context)
}

type InstructorHandlers struct {
	services appservices.Services
}

func NewInstructorHandlers(services appservices.Services) *InstructorHandlers {
	return &InstructorHandlers{services: services}
}

func (r *InstructorHandlers) InstructorRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

	instructorGroup := rg.Group("/instructor").Use(m.AuthJwtTokenMiddleware())
	{
		instructorGroup.POST("", m.RBACPermission("instructors:create"), r.CreateInstructorHandler)
		instructorGroup.GET("", m.RBACPermission("instructors:list"), r.GetInstructorsHandler)
		instructorGroup.GET("/summary", m.RBACPermission("instructors:list"), r.GetSummaryHandler)
		instructorGroup.GET("/:id", r.GetInstructorByIDHandler)
		instructorGroup.PUT("/:id", m.RBACPermission("instructors:update"), r.UpdateInstructorHandler)
		instructorGroup.DELETE("/:id", m.RBACPermission("instructors:delete"), r.DeleteInstructorHandler)
		instructorGroup.POST("/:id/specialty", m.RBACPermission("instructors:update"), r.AssignInstructorSpecialtyHandler)
		instructorGroup.DELETE("/:id/specialty/:specialty_id", m.RBACPermission("instructors:update"), r.DeleteInstructorSpecialtyHandler)
	}

	branchInstructorGroup := rg.Group("/branches/:id/instructor").Use(m.AuthJwtTokenMiddleware())
	{
		branchInstructorGroup.POST("", m.RBACPermission("instructors:create"), r.AssignInstructorsToBranchHandler)
		branchInstructorGroup.GET("", m.RBACPermission("instructors:list"), r.ListBranchInstructorsHandler)
		branchInstructorGroup.DELETE("/:instructor_id", m.RBACPermission("instructors:delete"), r.RemoveInstructorFromBranchHandler)
	}
}
