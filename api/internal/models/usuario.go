package models

import "time"

// Enumerations
type Permissao string

const (
	LEITURA         Permissao = "LEITURA"
	ESCRITA         Permissao = "ESCRITA"
	LEITURA_ESCRITA Permissao = "LEITURA_ESCRITA"
	SEM_PERMISSAO   Permissao = "SEM_PERMISSAO"
)

type TipoUsuario string

const (
	MORADOR TipoUsuario = "MORADOR"
	EXTERNO TipoUsuario = "EXTERNO"
)

type FaixaEtaria string

const (
	ADULTO  FaixaEtaria = "ADULTO"
	CRIANCA FaixaEtaria = "CRIANCA"
)

// Usuario struct
type Usuario struct {
	ID                   int
	Nome                 string
	Email                string
	Telefone             string
	DataNascimento       time.Time
	TipoUsuario          TipoUsuario
	FaixaEtaria          FaixaEtaria
	IsSindico            bool
	IsConselheiro        bool
	IsRepresentanteLegal bool
}

// Métodos para Usuario
func (u *Usuario) VerificarPermissoes() Permissao {
	if u.IsSindico {
		return LEITURA_ESCRITA
	} else if u.IsConselheiro {
		return LEITURA
	} else {
		return SEM_PERMISSAO
	}
}

func (u *Usuario) ListarApartamentosAssociados() []Apartamento {
	// Implementação fictícia para retornar apartamentos associados
	return []Apartamento{}
}
