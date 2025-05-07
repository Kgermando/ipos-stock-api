package models

import (
	"time"

	"gorm.io/gorm"
)

type Restitution struct {
	gorm.Model

	UUID            string      `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID         string      `json:"pos_uuid"`
	Pos             Pos         `gorm:"foreignKey:PosUUID"`
	ProductUUID     string      `json:"product_uuid"`
	Product         Product     `gorm:"foreignKey:ProductUUID"`
	Description     string      `json:"description"`
	Quantity        uint64      `gorm:"not null" json:"quantity"`
	FournisseurUUID string      `json:"fournisseur_uuid"`
	Fournisseur     Fournisseur `gorm:"foreignKey:FournisseurUUID"`
	PrixAchat       float64     `gorm:"not null" json:"prix_achat"`
	DateExpiration  time.Time   `gorm:"not null" json:"date_expiration"`
	Signature       string      `json:"signature"`
	CodeEntreprise  uint64      `json:"code_entreprise"`
}
