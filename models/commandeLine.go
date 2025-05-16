package models

import (
	"gorm.io/gorm"
)

type CommandeLine struct {
	gorm.Model

	UUID           string   `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	CommandeUUID   string   `gorm:"type:varchar(255);not null" json:"commande_uuid"`
	Commande       Commande `gorm:"foreignKey:CommandeUUID;references:UUID"` // Commande associée
	ProductUUID    string   `gorm:"type:varchar(255);not null" json:"product_uuid"`
	Product        Product  `gorm:"foreignKey:ProductUUID;references:UUID"` // Produit associé
	Quantity       uint64   `gorm:"not null" json:"quantity"`
	EntrepriseUUID string   `json:"entreprise_uuid"`
	Sync           bool     `gorm:"default:false" json:"sync"`
}
