package models

type Building struct {
	BaseModel
	Name string `json:"name" gorm:"not null;size:100" validate:"required,max=100"`

	// Chave estrangeira
	Apartments []Apartment `json:"apartments,omitempty" gorm:"foreignKey:BuildingID"`
}
