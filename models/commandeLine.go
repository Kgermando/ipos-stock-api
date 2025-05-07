package models

import (
	"gorm.io/gorm"
)

type CommandeLine struct {
	gorm.Model

	UUID         string   `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	CommandeUUID string   `json:"commande_uuid"`
	Commande     Commande `gorm:"foreignKey:CommandeUUID"`
	ProductUUID    string  `json:"product_uuid"`
	Product        Product `gorm:"foreignKey:ProductUUID"`
	Quantity       uint64  `gorm:"not null" json:"quantity"`
	CodeEntreprise uint64  `json:"code_entreprise"`
}
