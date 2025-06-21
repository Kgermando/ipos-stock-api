package models

import (
	"time"

	"gorm.io/gorm"
)

type StockEndommage struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PosUUID   string         `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos       Pos            `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	// Reference      uint64  `gorm:"not null" json:"reference"`          // Numero de reference du ravitaillement pour retrouver dans quel revitaillement le produit est endommagE
	ProductUUID    string  `gorm:"type:varchar(255);not null" json:"product_uuid"`
	Product        Product `gorm:"foreignKey:ProductUUID;references:UUID"` // Produit associé
	Quantity       float64 `gorm:"not null" json:"quantity"`
	PrixAchat      float64 `gorm:"not null" json:"prix_achat"`
	Raison         string  `json:"raison"` // Raison de l'endommagement
	Signature      string  `json:"signature"`
	EntrepriseUUID string  `json:"entreprise_uuid"`
	Sync           bool    `gorm:"default:false" json:"sync"`
}
