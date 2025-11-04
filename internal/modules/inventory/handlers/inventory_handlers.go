package inventoryhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
)

//
// EQUIPMENT CATALOG HANDLERS (Super Admin)
//

// @Summary		Create equipment type
// @Description	Adds a new equipment type to the global catalog (Super Admin only).
// @Tags			Equipment
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			equipment	body		CreateEquipmentPayload	true	"Equipment data"
// @Success		201			{object}	object{message=string}
// @Failure		400			{object}	object{error=string}
// @Failure		500			{object}	object{error=string}
// @Router			/equipment-types [post]
func (h *InventoryHandlers) CreateEquipmentHandler(c *gin.Context) {
	var payload CreateEquipmentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	equipment, err := payload.ToEquipment()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.EquipmentServices.CreateEquipment(c, equipment); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "equipment type created"})
}

// @Summary		Get equipment types
// @Description	Lists all available equipment types in the global catalog.
// @Tags			Equipment
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{data=[]inventorydomain.Equipment}
// @Failure		500	{object}	object{error=string}
// @Router			/equipment-types [get]
func (h *InventoryHandlers) GetEquipmentsHandler(c *gin.Context) {
	equipments, err := h.services.EquipmentServices.GetEquipments(c)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": equipments.Equipments})
}

// @Summary		Get equipment type by ID
// @Description	Retrieves a single equipment type from the global catalog by its UUID.
// @Tags			Equipment
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string	true	"Equipment UUID"
// @Success		200	{object}	object{data=EquipmentListResponse}
// @Failure		400	{object}	object{error=string}	"Error: Bad Request (e.g., invalid UUID)"
// @Failure		404	{object}	object{error=string}	"Error: Not Found"
// @Failure		500	{object}	object{error=string}	"Error: Internal Server Error"
// @Router			/equipment-types/{id} [get]
func (h *InventoryHandlers) GetEquipmentByID(c *gin.Context) {
	idStr := c.Param("id")
	equipmentID, err := uuid.Parse(idStr)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	equipment, err := h.services.EquipmentServices.GetEquipmentByID(c, equipmentID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	resp := EquipmentListResponse{
		EquipmentID: equipment.EquipmentID,
		Name:        equipment.Name,
		Category:    string(equipment.Category),
		Description: equipment.Description,
		Brand:       equipment.Brand,
		Model:       equipment.Model,
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// @Summary		Get branch inventory item by ID
// @Description	Retrieves a single inventory item from a branch by its UUID.
// @Tags			Branch Inventory
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path		string	true	"Branch UUID"
// @Param			inventoryId	path		string	true	"Inventory Item UUID"
// @Success		200			{object}	object{data=BranchInventoryResponse}
// @Failure		400			{object}	object{error=string}	"Error: Bad Request (e.g., invalid UUID)"
// @Failure		404			{object}	object{error=string}	"Error: Not Found"
// @Failure		500			{object}	object{error=string}	"Error: Internal Server Error"
// @Router			/branches/{id}/equipment/{inventoryId} [get]
func (h *InventoryHandlers) GetInventoryByID(c *gin.Context) {
	idStr := c.Param("inventoryId")
	inventoryID, err := uuid.Parse(idStr)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	equipment, err := h.services.InventoryServices.GetInventoryItemByID(c, inventoryID)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	var acquisitionDate, lastMaintenanceDate string
	if equipment.AcquisitionDate != nil {
		acquisitionDate = equipment.AcquisitionDate.Format("2006-01-02")
	}
	if equipment.LastMaintenanceDate != nil {
		lastMaintenanceDate = equipment.LastMaintenanceDate.Format("2006-01-02")
	}

	resp := BranchInventoryResponse{
		InventoryID:         equipment.InventoryID,
		EquipmentID:         equipment.EquipmentID,
		SerialNumber:        equipment.SerialNumber,
		Status:              string(equipment.Status),
		AcquisitionDate:     acquisitionDate,
		LastMaintenanceDate: lastMaintenanceDate,
		Notes:               equipment.Notes,
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// @Summary		Update equipment type
// @Description	Updates an existing equipment type in the global catalog.
// @Tags			Equipment
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id			path		string					true	"Equipment UUID"
// @Param			equipment	body		UpdateEquipmentPayload	true	"Equipment update payload"
// @Success		200			{object}	object{message=string}
// @Failure		400			{object}	object{error=string}
// @Failure		404			{object}	object{error=string}
// @Failure		500			{object}	object{error=string}
// @Router			/equipment-types/{id} [put]
func (h *InventoryHandlers) UpdateEquipmentHandler(c *gin.Context) {
	idStr := c.Param("id")
	equipmentID, err := uuid.Parse(idStr)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	var payload UpdateEquipmentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	equipment, err := payload.ToEquipment()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	equipment.EquipmentID = equipmentID

	if err := h.services.EquipmentServices.UpdateEquipment(c, equipment); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "equipment updated successfully"})
}

// @Summary		Delete equipment type
// @Description	Deletes an equipment type from the global catalog.
// @Tags			Equipment
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string	true	"Equipment UUID"
// @Success		204	{object}	nil
// @Failure		400	{object}	object{error=string}
// @Failure		404	{object}	object{error=string}
// @Failure		500	{object}	object{error=string}
// @Router			/equipment-types/{id} [delete]
func (h *InventoryHandlers) DeleteEquipmentHandler(c *gin.Context) {
	idStr := c.Param("id")
	equipmentID, err := uuid.Parse(idStr)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.EquipmentServices.DeleteEquipment(c, equipmentID); err != nil {
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

//
// BRANCH INVENTORY HANDLERS (Branch Admin)
//

// @Summary		Add item to branch inventory
// @Description	Adds a new equipment instance to a specific branch's inventory.
// @Tags			Branch Inventory
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id		path		string						true	"Branch UUID"
// @Param			item	body		CreateInventoryItemPayload	true	"Inventory item data"
// @Success		201		{object}	object{message=string}
// @Failure		400		{object}	object{error=string}
// @Failure		500		{object}	object{error=string}
// @Router			/branches/{id}/equipment [post]
func (h *InventoryHandlers) AddInventoryItemHandler(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	var payload CreateInventoryItemPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	item, err := payload.ToBranchInventory(branchID)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.InventoryServices.AddInventoryItem(c, item); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "inventory item added"})
}

// @Summary		List branch inventory
// @Description	Lists all equipment items for a given branch.
// @Tags			Branch Inventory
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string	true	"Branch UUID"
// @Success		200	{object}	object{data=[]inventorydomain.BranchInventory}
// @Failure		500	{object}	object{error=string}
// @Router			/branches/{id}/equipment [get]
func (h *InventoryHandlers) ListBranchInventoryHandler(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	inventory, err := h.services.InventoryServices.ListBranchInventory(c, branchID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventory.Inventory})
}

// @Summary		Update branch inventory item
// @Description	Updates an inventory item's status or details.
// @Tags			Branch Inventory
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id			path		string						true	"Branch UUID"
// @Param			inventoryId	path		string						true	"Inventory UUID"
// @Param			payload		body		UpdateInventoryItemPayload	true	"Update payload"
// @Success		200			{object}	object{message=string}
// @Failure		400			{object}	object{error=string}
// @Failure		404			{object}	object{error=string}
// @Failure		500			{object}	object{error=string}
// @Router			/branches/{id}/equipment/{inventoryId} [patch]
func (h *InventoryHandlers) UpdateInventoryItemHandler(c *gin.Context) {
	inventoryID, err := uuid.Parse(c.Param("inventoryId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	var payload UpdateInventoryItemPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	item, err := payload.ToBranchInventoryUpdate(inventoryID)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	item.BranchID = branchID
	if err := h.services.InventoryServices.UpdateInventoryItem(c, item); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "inventory item updated"})
}

// @Summary		Delete branch inventory item
// @Description	Deletes an equipment item from a branch's inventory.
// @Tags			Branch Inventory
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path		string	true	"Branch UUID"
// @Param			inventoryId	path		string	true	"Inventory UUID"
// @Success		204			{object}	nil
// @Failure		400			{object}	object{error=string}
// @Failure		404			{object}	object{error=string}
// @Failure		500			{object}	object{error=string}
// @Router			/branches/{id}/equipment/{inventoryId} [delete]
func (h *InventoryHandlers) DeleteInventoryItemHandler(c *gin.Context) {
	inventoryID, err := uuid.Parse(c.Param("inventoryId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.InventoryServices.DeleteInventoryItem(c, inventoryID, branchID); err != nil {
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
