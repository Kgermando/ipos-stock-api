package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	UUID           string  `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID        string  `json:"pos_uuid"`
	Pos            Pos     `gorm:"foreignKey:PosUUID"`
	Image          string  `json:"image"`
	Reference      string  `gorm:"not null" json:"reference"`
	Name           string  `gorm:"not null" json:"name"`
	Description    string  `gorm:"not null" json:"description"`
	UniteVente     string  `json:"unite_vente"`
	PrixVente      float64 `gorm:"not null" json:"prix_vente"`
	Tva            float64 `gorm:"default:0" json:"tva"`
	Stock          float64 `gorm:"default:0" json:"stock"`           // stock disponible
	StockEndommage float64 `gorm:"default:0" json:"stock_endommage"` // stock endommage
	Restitution    float64 `gorm:"default:0" json:"restitution"`     // stock restitution

	CodeEntreprise uint64 `json:"code_entreprise"`
	Signature      string `json:"signature"`
	Sync            bool       `gorm:"default:false" json:"sync"`

	Stocks          []Stock          `gorm:"foreignKey:ProductUUID"`
	StockEndommages []StockEndommage `gorm:"foreignKey:ProductUUID"`
	CommadeLines    []CommandeLine   `gorm:"foreignKey:ProductUUID"`
	Restitutions    []Restitution    `gorm:"foreignKey:ProductUUID"`
}
