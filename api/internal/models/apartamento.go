package models

// Apartamento struct
type Apartamento struct {
	ID                 int
	Numero             string
	Bloco              string
	Andar              int
	Usuarios           []Usuario
	Carros             []Carro
	Motos              []Moto
	Pets               []Pet
	Bicicletas         []Bicicleta
	RepresentanteLegal *Usuario
}

// Métodos para Apartamento
func (a *Apartamento) AdicionarUsuario(usuario Usuario) {
	a.Usuarios = append(a.Usuarios, usuario)
}

func (a *Apartamento) RemoverUsuario(usuario Usuario) {
	for i, u := range a.Usuarios {
		if u.ID == usuario.ID {
			a.Usuarios = append(a.Usuarios[:i], a.Usuarios[i+1:]...)
			break
		}
	}
}

func (a *Apartamento) DefinirRepresentanteLegal(usuario Usuario) {
	a.RepresentanteLegal = &usuario
}

func (a *Apartamento) AdicionarVeiculo(veiculo interface{}) {
	switch v := veiculo.(type) {
	case Carro:
		a.Carros = append(a.Carros, v)
	case Moto:
		a.Motos = append(a.Motos, v)
	}
}

func (a *Apartamento) AdicionarPet(pet Pet) {
	a.Pets = append(a.Pets, pet)
}

func (a *Apartamento) AdicionarBicicleta(bicicleta Bicicleta) {
	a.Bicicletas = append(a.Bicicletas, bicicleta)
}
