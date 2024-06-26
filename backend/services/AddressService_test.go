package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/arthur-tragante/liven-code-test/models"
	"github.com/arthur-tragante/liven-code-test/services"
	"github.com/arthur-tragante/liven-code-test/testutils"
)

type ServiceTestSuite struct {
	suite.Suite
	TestDBSetup    *testutils.TestDBSetup
	UserService    *services.UserService
	AddressService *services.AddressService
	DB             *gorm.DB
}

func (suite *ServiceTestSuite) SetupSuite() {
	suite.TestDBSetup = testutils.SetupTestDB(assert.New(suite.T()))
	suite.DB = suite.TestDBSetup.DB
	suite.UserService = &services.UserService{
		DB:        suite.DB,
		JWTSecret: "testsecret",
	}
	suite.AddressService = &services.AddressService{
		DB: suite.DB,
	}
}

func (suite *ServiceTestSuite) TearDownSuite() {
	suite.TestDBSetup.TearDown(assert.New(suite.T()))
}

func (suite *ServiceTestSuite) SetupTest() {
	err := suite.DB.Exec("DELETE FROM addresses").Error
	assert.NoError(suite.T(), err)
	err = suite.DB.Exec("DELETE FROM users").Error
	assert.NoError(suite.T(), err)
}

func (suite *ServiceTestSuite) TestCreateAddress() {
	user := &models.User{
		Name:     "Test User",
		Email:    "test.user+create@example.com",
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

	var createdAddress models.Address
	err = suite.DB.First(&createdAddress, "user_id = ?", user.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), address.Street, createdAddress.Street)
	assert.Equal(suite.T(), address.City, createdAddress.City)
}

func (suite *ServiceTestSuite) TestGetAddressByID() {
	user := &models.User{
		Name:     "Test User",
		Email:    "test.user+getbyid@example.com",
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

	fetchedAddress, err := suite.AddressService.GetAddressByID(address.AddressID, user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), address.Street, fetchedAddress.Street)
	assert.Equal(suite.T(), address.City, fetchedAddress.City)
}

func (suite *ServiceTestSuite) TestGetAllAddresses() {
	user := &models.User{
		Name:     "Test User",
		Email:    "test.user+getall@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	address1 := &models.Address{
		UserID:     user.ID,
		Street:     "123 Test St",
		Number:     "1",
		Complement: "Apt 1",
		City:       "Test City",
		State:      "Test State",
		Zipcode:    "12345",
		Country:    "Test Country",
	}

	address2 := &models.Address{
		UserID:     user.ID,
		Street:     "456 Another St",
		Number:     "2",
		Complement: "Apt 2",
		City:       "Another City",
		State:      "Another State",
		Zipcode:    "67890",
		Country:    "Another Country",
	}

	err = suite.AddressService.CreateAddress(address1)
	assert.NoError(suite.T(), err)

	err = suite.AddressService.CreateAddress(address2)
	assert.NoError(suite.T(), err)

	addresses, err := suite.AddressService.GetAllAddresses(user.ID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), addresses, 2)
}

func (suite *ServiceTestSuite) TestUpdateAddress() {
	user := &models.User{
		Name:     "Test User",
		Email:    "test.user+update@example.com",
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

	err = suite.AddressService.UpdateAddress(address.AddressID, user.ID, updatedData)
	assert.NoError(suite.T(), err)

	var updatedAddress models.Address
	err = suite.DB.First(&updatedAddress, "address_id = ?", address.AddressID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedData.Street, updatedAddress.Street)
	assert.Equal(suite.T(), updatedData.City, updatedAddress.City)
}

func (suite *ServiceTestSuite) TestDeleteAddress() {
	user := &models.User{
		Name:     "Test User",
		Email:    "test.user+delete@example.com",
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

	err = suite.AddressService.DeleteAddress(address.AddressID, user.ID)
	assert.NoError(suite.T(), err)

	var deletedAddress models.Address
	err = suite.DB.First(&deletedAddress, "address_id = ?", address.AddressID).Error
	assert.Error(suite.T(), err)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
