package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	UUID        string  `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID     string  `json:"pos_uuid"`
	Pos         Pos     `gorm:"foreignKey:PosUUID"`
	Image       string  `json:"image"`
	Reference   string  `gorm:"not null" json:"reference"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	UniteVente  string  `json:"unite_vente"`
	PrixVente   float64 `gorm:"not null" json:"prix_vente"`
	Tva         float64 `gorm:"default:0" json:"tva"`

	CodeEntreprise uint64 `json:"code_entreprise"`
	Signature      string `json:"signature"`

	Stocks          []Stock          `gorm:"foreignKey:ProductUUID"`
	StockEndommages []StockEndommage `gorm:"foreignKey:ProductUUID"`
	CommadeLines    []CommandeLine   `gorm:"foreignKey:ProductUUID"`
	Restitutions    []Restitution    `gorm:"foreignKey:ProductUUID"`
}
