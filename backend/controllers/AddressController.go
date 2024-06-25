package controllers

import (
	"net/http"
	"strconv"

	"github.com/arthur-tragante/liven-code-test/models"
	"github.com/arthur-tragante/liven-code-test/services"
	"github.com/gin-gonic/gin"
)

type AddressController struct {
	AddressService *services.AddressService
}

func (ctrl *AddressController) CreateAddress(c *gin.Context) {
	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)
	address.UserID = userID

	if err := ctrl.AddressService.CreateAddress(&address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, address)
}

func (ctrl *AddressController) GetAddress(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	addressIDStr := c.Param("id")
	if addressIDStr != "" {
		addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
			return
		}

		address, err := ctrl.AddressService.GetAddressByID(uint(addressID), userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusOK, address)
		return
	}

	addresses, err := ctrl.AddressService.GetAllAddresses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (ctrl *AddressController) UpdateAddress(c *gin.Context) {
	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	var updatedData models.Address
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	if err := ctrl.AddressService.UpdateAddress(uint(addressID), userID, &updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedData)
}

func (ctrl *AddressController) DeleteAddress(c *gin.Context) {
	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	userID := c.MustGet("userID").(uint)

	if err := ctrl.AddressService.DeleteAddress(uint(addressID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}