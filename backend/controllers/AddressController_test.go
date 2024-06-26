package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type AddressControllerTestSuite struct {
	suite.Suite
	TestDBSetup       *testutils.TestDBSetup
	UserService       *services.UserService
	AddressService    *services.AddressService
	AddressController *controllers.AddressController
	DB                *gorm.DB
}

func (suite *AddressControllerTestSuite) SetupSuite() {
	suite.TestDBSetup = testutils.SetupTestDB(assert.New(suite.T()))
	suite.DB = suite.TestDBSetup.DB
	suite.UserService = &services.UserService{
		DB:        suite.DB,
		JWTSecret: "testsecret",
	}
	suite.AddressService = &services.AddressService{
		DB: suite.DB,
	}
	suite.AddressController = &controllers.AddressController{
		AddressService: suite.AddressService,
	}
}

func (suite *AddressControllerTestSuite) TearDownSuite() {
	suite.TestDBSetup.TearDown(assert.New(suite.T()))
}
func (suite *AddressControllerTestSuite) SetupTest() {
	err := suite.DB.Exec("DELETE FROM addresses").Error
	assert.NoError(suite.T(), err)
	err = suite.DB.Exec("DELETE FROM users").Error
	assert.NoError(suite.T(), err)
}

func (suite *AddressControllerTestSuite) TestCreateAddress_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	address := &models.Address{
		Street:     "123 Test St",
		Number:     "1",
		Complement: "Apt 1",
		City:       "Test City",
		State:      "Test State",
		Zipcode:    "12345",
		Country:    "Test Country",
	}

	jsonValue, _ := json.Marshal(address)
	c.Request, _ = http.NewRequest("POST", "/address", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", user.ID)

	suite.AddressController.CreateAddress(c)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdAddress models.Address
	err = json.Unmarshal(w.Body.Bytes(), &createdAddress)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), address.Street, createdAddress.Street)
	assert.Equal(suite.T(), user.ID, createdAddress.UserID)
}

func (suite *AddressControllerTestSuite) TestCreateAddress_InvalidJSON() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/address", bytes.NewBufferString("{invalid json}"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", uint(1))

	suite.AddressController.CreateAddress(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AddressControllerTestSuite) TestGetAddress_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "Jane Doe",
		Email:    "jane.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	address := &models.Address{
		UserID:     user.ID,
		Street:     "123 Test St",
		Number:     "1",
		Complement: "Apt 1",
		City:       "Test City",
		State:      "Test State",
		Zipcode:    "12345",
		Country:    "Test Country",
	}

	err = suite.AddressService.CreateAddress(address)
	assert.NoError(suite.T(), err)

	c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/address/%d", address.AddressID), nil)
	c.Set("userID", user.ID)
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", address.AddressID)}}

	suite.AddressController.GetAddress(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var fetchedAddress models.Address
	err = json.Unmarshal(w.Body.Bytes(), &fetchedAddress)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), address.Street, fetchedAddress.Street)
}

func (suite *AddressControllerTestSuite) TestGetAddress_InvalidID() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/address/invalid", nil)
	c.Set("userID", uint(1))
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	suite.AddressController.GetAddress(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AddressControllerTestSuite) TestGetAddress_NotFound() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/address/999", nil)
	c.Set("userID", uint(1))
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	suite.AddressController.GetAddress(c)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *AddressControllerTestSuite) TestUpdateAddress_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "Jim Doe",
		Email:    "jim.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	address := &models.Address{
		UserID:     user.ID,
		Street:     "123 Test St",
		Number:     "1",
		Complement: "Apt 1",
		City:       "Test City",
		State:      "Test State",
		Zipcode:    "12345",
		Country:    "Test Country",
	}

	err = suite.AddressService.CreateAddress(address)
	assert.NoError(suite.T(), err)

	updatedData := &models.Address{
		Street:  "Updated Street",
		City:    "Updated City",
		State:   "Updated State",
		Zipcode: "54321",
		Country: "Updated Country",
	}

	jsonValue, _ := json.Marshal(updatedData)
	c.Request, _ = http.NewRequest("PUT", fmt.Sprintf("/address/%d", address.AddressID), bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", user.ID)
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", address.AddressID)}}

	suite.AddressController.UpdateAddress(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var updatedAddress models.Address
	err = json.Unmarshal(w.Body.Bytes(), &updatedAddress)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedData.Street, updatedAddress.Street)
	assert.Equal(suite.T(), updatedData.City, updatedAddress.City)
}

func (suite *AddressControllerTestSuite) TestUpdateAddress_InvalidID() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/address/invalid", bytes.NewBufferString("{}"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", uint(1))
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	suite.AddressController.UpdateAddress(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AddressControllerTestSuite) TestUpdateAddress_BadRequest() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/address/1", bytes.NewBufferString("{invalid json}"))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", uint(1))
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	suite.AddressController.UpdateAddress(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AddressControllerTestSuite) TestDeleteAddress_Success() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	user := &models.User{
		Name:     "Jake Doe",
		Email:    "jake.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	address := &models.Address{
		UserID:     user.ID,
		Street:     "123 Test St",
		Number:     "1",
		Complement: "Apt 1",
		City:       "Test City",
		State:      "Test State",
		Zipcode:    "12345",
		Country:    "Test Country",
	}

	err = suite.AddressService.CreateAddress(address)
	assert.NoError(suite.T(), err)

	c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/address/%d", address.AddressID), nil)
	c.Set("userID", user.ID)
	c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", address.AddressID)}}

	suite.AddressController.DeleteAddress(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *AddressControllerTestSuite) TestDeleteAddress_InvalidID() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("DELETE", "/address/invalid", nil)
	c.Set("userID", uint(1))
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	suite.AddressController.DeleteAddress(c)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *AddressControllerTestSuite) TestDeleteAddress_NotFound() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("DELETE", "/address/999", nil)
	c.Set("userID", uint(1))
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	suite.AddressController.DeleteAddress(c)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestAddressControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AddressControllerTestSuite))
}
