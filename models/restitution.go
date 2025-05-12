package models

import (
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
	PrixAchat       float64     `gorm:"not null" json:"prix_achat"`
	Motif           string      `json:"motif"`
	FournisseurUUID string      `json:"fournisseur_uuid"`
	Fournisseur     Fournisseur `gorm:"foreignKey:FournisseurUUID"`
	Signature       string      `json:"signature"`
	CodeEntreprise  uint64      `json:"code_entreprise"`
	Sync            bool       `gorm:"default:false" json:"sync"`
}
