package models

import "time"

type UserApartment struct {
	UserID      uint      `json:"user_id" gorm:"primaryKey"`
	ApartmentID uint      `json:"apartment_id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Apartment Apartment `json:"apartment,omitempty" gorm:"foreignKey:ApartmentID"`
}

func (UserApartment) TableName() string {
	return "user_apartments"
}
