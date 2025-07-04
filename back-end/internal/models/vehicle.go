package models

type Vehicle struct {
	ID            int
	Plate         string
	Model         string
	Color         string
	Year          int
	ParkingNumber string
}

type Car struct {
	Vehicle
}

type Motorcycle struct {
	Vehicle
}
