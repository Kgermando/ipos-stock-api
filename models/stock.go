package models

import (
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model

	UUID            string      `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID         string      `json:"pos_uuid"`
	Pos             Pos         `gorm:"foreignKey:PosUUID"`
	ProductUUID     string      `json:"product_uuid"`
	Product         Product     `gorm:"foreignKey:ProductUUID"`
	Description     string      `json:"description"`
	Quantity        float64     `gorm:"not null" json:"quantity"`
	PrixAchat       float64     `gorm:"not null" json:"prix_achat"`
	DateExpiration  time.Time   `gorm:"not null" json:"date_expiration"`
	FournisseurUUID string      `json:"fournisseur_uuid"`
	Fournisseur     Fournisseur `gorm:"foreignKey:FournisseurUUID"`
	Signature       string      `json:"signature"`
	CodeEntreprise  uint64      `json:"code_entreprise"`
	Sync            bool       `gorm:"default:false" json:"sync"`
}

type FournisseurStock struct {
	Name           string  `json:"name"`
	Telephone      string  `json:"telephone"`
	TypeFourniture string  `json:"type_fourniture"`
	TotalValue     float64 `json:"total_value"`
}
