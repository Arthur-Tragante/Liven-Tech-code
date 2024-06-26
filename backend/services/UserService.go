package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/arthur-tragante/liven-code-test/models"
)

type UserService struct {
	DB        *gorm.DB
	JWTSecret string
}

// function to handle the registrations of the user

func (s *UserService) Register(user *models.User) error {
	// Using bcrypt to hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Password hashing error:", err)
		return err
	}
	user.Password = string(hashedPassword)
	return s.DB.Create(user).Error
}

func (s *UserService) Login(email, password string) (string, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("Email lookup error:", err)
		return "", errors.New("invalid email or password")
	}

	// Comparing the password in database with the password received in the request
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Password comparison failed:", err)
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		fmt.Println("Token signing error:", err)
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("Addresses").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(userID uint, updatedData *models.User) error {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	user.Name = updatedData.Name
	user.Email = updatedData.Email

	if updatedData.Password != "" {
		// Same logic for previous password hashing
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedData.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.DB.Save(&user).Error
}

func (s *UserService) DeleteUser(userID uint) error {
	return s.DB.Where("id = ?", userID).Delete(&models.User{}).Error
}
