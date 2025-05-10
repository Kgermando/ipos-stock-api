package models

import (
	"gorm.io/gorm"
)

type StockEndommage struct {
	gorm.Model

	UUID           string  `gorm:"type:text;not null;unique" json:"uuid"`
	PosUUID        string  `json:"pos_uuid"`
	Pos            Pos     `gorm:"foreignKey:PosUUID"`
	ProductUUID    string  `json:"product_uuid"`
	Product        Product `gorm:"foreignKey:ProductUUID"`
	Quantity       float64 `gorm:"not null" json:"quantity"`
	PrixAchat      float64 `gorm:"not null" json:"prix_achat"`
	Raison         string  `json:"raison"` // Raison de l'endommagement
	Signature      string  `json:"signature"`
	CodeEntreprise uint64  `json:"code_entreprise"`
}
