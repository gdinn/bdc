package models

import (
	"time"

	"github.com/google/uuid"
)

type UserApartment struct {
	UserID      uuid.UUID `json:"user_id" gorm:"primaryKey"`
	ApartmentID uuid.UUID `json:"apartment_id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Apartment Apartment `json:"apartment,omitempty" gorm:"foreignKey:ApartmentID"`
}

func (UserApartment) TableName() string {
	return "user_apartments"
}
