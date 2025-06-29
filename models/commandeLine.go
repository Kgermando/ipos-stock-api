package models

import (
	"time"

	"gorm.io/gorm"
)

type CommandeLine struct {
	UUID           string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	CommandeUUID   string         `gorm:"type:varchar(255);not null" json:"commande_uuid"`
	Commande       Commande       `gorm:"foreignKey:CommandeUUID;references:UUID"` // Commande associée
	ProductUUID    string         `gorm:"type:varchar(255);not null" json:"product_uuid"`
	Product        Product        `gorm:"foreignKey:ProductUUID;references:UUID"` // Produit associé
	Quantity       uint64         `gorm:"not null" json:"quantity"`
	EntrepriseUUID string         `json:"entreprise_uuid"`
	PosUUID        string         `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos            Pos            `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Sync           bool           `gorm:"default:false" json:"sync"`
}
