package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditInfo struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	IsActive  bool       `json:"is_active"`
}

type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate hook para garantir que o UUID seja gerado
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// IsDeleted verifica se o registro foi excluído (soft delete)
func (b *BaseModel) IsDeleted() bool {
	return b.DeletedAt != nil
}

// SoftDelete marca o registro como excluído
func (b *BaseModel) SoftDelete() {
	now := time.Now()
	b.DeletedAt = &now
}

// Restore restaura um registro excluído
func (b *BaseModel) Restore() {
	b.DeletedAt = nil
}

// GetAuditInfo retorna informações de auditoria
func (b *BaseModel) GetAuditInfo() AuditInfo {
	return AuditInfo{
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		DeletedAt: b.DeletedAt,
		IsActive:  !b.IsDeleted(),
	}
}
