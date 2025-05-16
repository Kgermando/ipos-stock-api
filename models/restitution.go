package models

import (
	"gorm.io/gorm"
)

type Restitution struct {
	gorm.Model

	UUID            string      `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID         string      `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos             Pos         `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	ProductUUID     string      `gorm:"type:varchar(255);not null" json:"product_uuid"`
	Product         Product     `gorm:"foreignKey:ProductUUID;references:UUID"` // Produit associé
	Description     string      `json:"description"`
	Quantity        uint64      `gorm:"not null" json:"quantity"`
	PrixAchat       float64     `gorm:"not null" json:"prix_achat"`
	Motif           string      `json:"motif"`
	FournisseurUUID string      `gorm:"type:varchar(255);not null" json:"fournisseur_uuid"`
	Fournisseur     Fournisseur `gorm:"foreignKey:FournisseurUUID;references:UUID"` // Fournisseur associé
	Signature       string      `json:"signature"`
	EntrepriseUUID  string      `json:"entreprise_uuid"`
	Sync            bool        `gorm:"default:false" json:"sync"`
}
