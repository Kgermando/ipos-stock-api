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
	Quantity        float64      `gorm:"not null" json:"quantity"`
	FournisseurUUID string      `json:"fournisseur_uuid"`
	Fournisseur     Fournisseur `gorm:"foreignKey:FournisseurUUID"`
	PrixAchat       float64     `gorm:"not null" json:"prix_achat"`
	DateExpiration  time.Time   `gorm:"not null" json:"date_expiration"`
	Endomage        bool        `json:"endomage"`
	Signature       string      `json:"signature"`
	CodeEntreprise  uint64      `json:"code_entreprise"`
}

type FournisseurStock struct {
	Name           string  `json:"name"`
	Telephone      string  `json:"telephone"`
	TypeFourniture string  `json:"type_fourniture"`
	TotalValue     float64 `json:"total_value"`
}
