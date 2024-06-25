package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	AddressID  uint           `gorm:"primaryKey" json:"address_id"`
	UserID     uint           `json:"user_id"`
	Street     string         `json:"street"`
	Number     string         `json:"number"`
	Complement string         `json:"complement"`
	City       string         `json:"city"`
	State      string         `json:"state"`
	Zipcode    string         `json:"zipcode"`
	Country    string         `json:"country"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}