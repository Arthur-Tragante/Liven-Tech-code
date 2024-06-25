package controllers

import (
	"fmt"
	"net/http"

	"github.com/arthur-tragante/liven-code-test/models"
	"github.com/arthur-tragante/liven-code-test/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func (ctrl *UserController) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("RegisterUser received: %+v\n", user)

	if err := ctrl.UserService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("LoginUser received: %+v\n", loginData)

	token, err := ctrl.UserService.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	user, err := ctrl.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	var updatedData models.User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	if err := ctrl.UserService.UpdateUser(userID, &updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedData)
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	if err := ctrl.UserService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
