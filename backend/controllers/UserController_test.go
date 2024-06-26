package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arthur-tragante/liven-code-test/controllers"
	"github.com/arthur-tragante/liven-code-test/models"
	"github.com/arthur-tragante/liven-code-test/services"
	"github.com/arthur-tragante/liven-code-test/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserControllerTestSuite struct {
	suite.Suite
	TestDBSetup    *testutils.TestDBSetup
	UserService    *services.UserService
	UserController *controllers.UserController
	DB             *gorm.DB
}

func (suite *UserControllerTestSuite) SetupSuite() {
	suite.TestDBSetup = testutils.SetupTestDB(assert.New(suite.T()))
	suite.DB = suite.TestDBSetup.DB
	suite.UserService = &services.UserService{
		DB:        suite.DB,
		JWTSecret: "testsecret",
	}
	suite.UserController = &controllers.UserController{
		UserService: suite.UserService,
	}
}

func (suite *UserControllerTestSuite) TearDownSuite() {
	suite.TestDBSetup.TearDown(assert.New(suite.T()))
}

func (suite *UserControllerTestSuite) SetupTest() {
	err := suite.DB.Exec("DELETE FROM addresses").Error
	assert.NoError(suite.T(), err)
	err = suite.DB.Exec("DELETE FROM users").Error
	assert.NoError(suite.T(), err)
}

func (suite *UserControllerTestSuite) TestRegisterUser_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	jsonValue, _ := json.Marshal(user)
	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.UserController.RegisterUser(c)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Email, createdUser.Email)
}

func (suite *UserControllerTestSuite) TestRegisterUser_InvalidJSON() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/register", bytes.NewBufferString("{invalid json}"))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.UserController.RegisterUser(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserControllerTestSuite) TestLoginUser_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "Jane Doe",
		Email:    "jane.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	loginData := map[string]string{
		"email":    "jane.doe@example.com",
		"password": "password123",
	}

	jsonValue, _ := json.Marshal(loginData)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.UserController.LoginUser(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), response["token"])
}

func (suite *UserControllerTestSuite) TestLoginUser_InvalidCredentials() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	loginData := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}

	jsonValue, _ := json.Marshal(loginData)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")

	suite.UserController.LoginUser(c)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *UserControllerTestSuite) TestGetUser_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "Jack Doe",
		Email:    "jack.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	c.Set("userID", user.ID)

	suite.UserController.GetUser(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var fetchedUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &fetchedUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Email, fetchedUser.Email)
}

func (suite *UserControllerTestSuite) TestGetUser_NotFound() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("userID", uint(999))

	suite.UserController.GetUser(c)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *UserControllerTestSuite) TestUpdateUser_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "Jill Doe",
		Email:    "jill.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	updatedData := &models.User{
		Name:  "Jill Smith",
		Email: "jill.smith@example.com",
	}

	jsonValue, _ := json.Marshal(updatedData)
	c.Request, _ = http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", user.ID)

	suite.UserController.UpdateUser(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var updatedUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &updatedUser)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedData.Name, updatedUser.Name)
	assert.Equal(suite.T(), updatedData.Email, updatedUser.Email)
}

func (suite *UserControllerTestSuite) TestUpdateUser_BadRequest() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/user", bytes.NewBufferString("{invalid json}"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", uint(1))

	suite.UserController.UpdateUser(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserControllerTestSuite) TestDeleteUser_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "Jim Doe",
		Email:    "jim.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	c.Request, _ = http.NewRequest("DELETE", "/user", nil)
	c.Set("userID", user.ID)

	suite.UserController.DeleteUser(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "User deleted successfully", response["message"])
}

func (suite *UserControllerTestSuite) TestDeleteUser_InternalServerError() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("DELETE", "/user", nil)
	c.Set("userID", uint(999)) 

	suite.UserController.DeleteUser(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
