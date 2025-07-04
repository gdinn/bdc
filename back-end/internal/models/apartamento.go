package models

type Apartment struct {
	ID          int
	Number      string
	Building    string
	Users       []User
	Cars        []Car
	Motorcycles []Motorcycle
	Pets        []Pet
	Bicycles    []Bicycle
	LegalRep    *User
}

func (a *Apartment) AddUser(user User) {
	a.Users = append(a.Users, user)
}

func (a *Apartment) RemoverUser(user User) {
	for i, u := range a.Users {
		if u.ID == user.ID {
			a.Users = append(a.Users[:i], a.Users[i+1:]...)
			break
		}
	}
}

func (a *Apartment) SetLegalRep(user User) {
	a.LegalRep = &user
}

func (a *Apartment) AddVehicle(vehicle interface{}) {
	switch v := vehicle.(type) {
	case Car:
		a.Cars = append(a.Cars, v)
	case Motorcycle:
		a.Motorcycles = append(a.Motorcycles, v)
	}
}

func (a *Apartment) AddPet(pet Pet) {
	a.Pets = append(a.Pets, pet)
}

func (a *Apartment) AddBicycle(bicycle Bicycle) {
	a.Bicycles = append(a.Bicycles, bicycle)
}
