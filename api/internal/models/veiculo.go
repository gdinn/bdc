package models

// Veiculo struct base
type Veiculo struct {
	ID     int
	Placa  string
	Modelo string
	Cor    string
	Ano    int
}

// Carro struct
type Carro struct {
	Veiculo
	NumeroVaga string
}

// Moto struct
type Moto struct {
	Veiculo
	NumeroVaga string
}
