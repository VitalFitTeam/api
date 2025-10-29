package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
)

// @Summary		List all roles in the system
// @Description	List all roles in the system
// @Tags			Admin
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Success		200	{object}	object{data=[]BranchAdminResponse}	"succed response"
// @Failure		500	{object}	object{error=string}				"Error: internal server error"
// @Router			/admin/roles [get]
func (r *AuthHandlers) GetRolesHandler(c *gin.Context) {
	roles, err := r.services.UserServices.GetRoles(c)
	if err != nil {
		r.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(200, gin.H{
		"data": roles,
	})
}

func (r *AuthHandlers) CreateRoleHandler(c *gin.Context) {
	var payload CreateRolesPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	role, err := payload.createRole()
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := r.services.UserServices.CreateRole(ctx, role); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			r.services.LogErrors.BadRequestResponse(c, err)
			return
		case shared_errors.ErrConflict:
			r.services.LogErrors.ConflictResponse(c, err)
			return
		default:
			r.services.LogErrors.InternalServerError(c, err)
			return
		}
	}
	c.JSON(http.StatusCreated, nil)
}

func (r *AuthHandlers) GetRoleByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	role, err := r.services.UserServices.GetRoleByID(ctx, roleID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			r.services.LogErrors.NotFoundResponse(c)
			return
		default:
			r.services.LogErrors.InternalServerError(c, err)
			return
		}
	}
	c.JSON(http.StatusOK, role)
}

func (r *AuthHandlers) UpdateRoleHandler(c *gin.Context) {
	ctx := c.Request.Context()
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	var payload CreateRolesPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	role, err := payload.createRole()
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	role.RoleID = roleID

	err = r.services.UserServices.UpdateRole(ctx, role)
	if err != nil {
		r.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)

}

func (r *AuthHandlers) DeleteRoleHandler(c *gin.Context) {
	ctx := c.Request.Context()
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	if err := r.services.UserServices.DeleteRole(ctx, roleID); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			r.services.LogErrors.NotFoundResponse(c)
			return
		default:
			r.services.LogErrors.InternalServerError(c, err)
			return
		}

	}

}

func (r *AuthHandlers) GetPermissionsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	permissions, err := r.services.UserServices.GetPermissions(ctx)
	if err != nil {
		r.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": permissions,
	})
}

func (r *AuthHandlers) AssignRolePermissionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	var payload PermissionsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	permissions, err := payload.toPermission()
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := r.services.UserServices.AssignRolePermission(ctx, roleID, permissions); err != nil {
		r.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (r *AuthHandlers) DeleteRolePermissionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	var payload PermissionsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	permissions, err := payload.toPermission()
	if err != nil {
		r.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := r.services.UserServices.DeleteRolePermission(ctx, roleID, permissions); err != nil {
		r.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
