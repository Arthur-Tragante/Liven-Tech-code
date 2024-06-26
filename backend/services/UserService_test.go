package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/arthur-tragante/liven-code-test/models"
	"github.com/arthur-tragante/liven-code-test/services"
	"github.com/arthur-tragante/liven-code-test/testutils"
)

type UserServiceTestSuite struct {
	suite.Suite
	TestDBSetup *testutils.TestDBSetup
	UserService *services.UserService
	DB          *gorm.DB
}

func (suite *UserServiceTestSuite) SetupSuite() {
	suite.TestDBSetup = testutils.SetupTestDB(assert.New(suite.T()))
	suite.DB = suite.TestDBSetup.DB
	suite.UserService = &services.UserService{
		DB:        suite.DB,
		JWTSecret: "testsecret",
	}
}

func (suite *UserServiceTestSuite) TearDownSuite() {
	suite.TestDBSetup.TearDown(assert.New(suite.T()))
}

func (suite *UserServiceTestSuite) TestRegister() {
	user := &models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	var createdUser models.User
	err = suite.DB.First(&createdUser, "email = ?", "john.doe@example.com").Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Name, createdUser.Name)
	assert.Equal(suite.T(), user.Email, createdUser.Email)
	assert.Equal(suite.T(), user.Password, createdUser.Password)

	err = bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte("password123"))
	assert.NoError(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestLogin() {
	user := &models.User{
		Name:     "Jane Doe",
		Email:    "jane.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	token, err := suite.UserService.Login("jane.doe@example.com", "password123")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)

	token, err = suite.UserService.Login("jane.doe@example.com", "wrongpassword")
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func (suite *UserServiceTestSuite) TestGetUserByID() {
	user := &models.User{
		Name:     "Jack Doe",
		Email:    "jack.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	var createdUser models.User
	err = suite.DB.First(&createdUser, "email = ?", "jack.doe@example.com").Error
	assert.NoError(suite.T(), err)

	fetchedUser, err := suite.UserService.GetUserByID(createdUser.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdUser.Email, fetchedUser.Email)
	assert.Equal(suite.T(), createdUser.Name, fetchedUser.Name)

	fetchedUser, err = suite.UserService.GetUserByID(9999)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), fetchedUser)
}

func (suite *UserServiceTestSuite) TestUpdateUser() {
	user := &models.User{
		Name:     "Jill Doe",
		Email:    "jill.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	var createdUser models.User
	err = suite.DB.First(&createdUser, "email = ?", "jill.doe@example.com").Error
	assert.NoError(suite.T(), err)

	updatedData := &models.User{
		Name:     "Jill Smith",
		Email:    "jill.smith@example.com",
		Password: "newpassword",
	}

	err = suite.UserService.UpdateUser(createdUser.ID, updatedData)
	assert.NoError(suite.T(), err)

	var updatedUser models.User
	err = suite.DB.First(&updatedUser, "id = ?", createdUser.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedData.Name, updatedUser.Name)
	assert.Equal(suite.T(), updatedData.Email, updatedUser.Email)

	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword"))
	assert.NoError(suite.T(), err)
}

func (suite *UserServiceTestSuite) TestDeleteUser() {
	user := &models.User{
		Name:     "Jim Doe",
		Email:    "jim.doe@example.com",
		Password: "password123",
	}

	err := suite.UserService.Register(user)
	assert.NoError(suite.T(), err)

	var createdUser models.User
	err = suite.DB.First(&createdUser, "email = ?", "jim.doe@example.com").Error
	assert.NoError(suite.T(), err)

	err = suite.UserService.DeleteUser(createdUser.ID)
	assert.NoError(suite.T(), err)

	var deletedUser models.User
	err = suite.DB.First(&deletedUser, "id = ?", createdUser.ID).Error
	assert.Error(suite.T(), err)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
