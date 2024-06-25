package services

import (
	"gorm.io/gorm"

	"github.com/arthur-tragante/liven-code-test/models"
)

type AddressService struct {
	DB *gorm.DB
}

func (s *AddressService) CreateAddress(address *models.Address) error {
	return s.DB.Create(address).Error
}

func (s *AddressService) GetAddressByID(addressID, userID uint) (*models.Address, error) {
	var address models.Address
	if err := s.DB.Where("address_id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (s *AddressService) GetAllAddresses(userID uint) ([]models.Address, error) {
	var addresses []models.Address
	if err := s.DB.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (s *AddressService) UpdateAddress(addressID, userID uint, updatedData *models.Address) error {
	return s.DB.Model(&models.Address{}).Where("address_id = ? AND user_id = ?", addressID, userID).Updates(updatedData).Error
}

func (s *AddressService) DeleteAddress(addressID, userID uint) error {
	return s.DB.Where("address_id = ? AND user_id = ?", addressID, userID).Delete(&models.Address{}).Error
}

func (s *AddressService) GetUserWithAddresses(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("Addresses").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}