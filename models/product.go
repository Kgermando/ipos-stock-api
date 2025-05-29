package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	UUID              string  `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID           string  `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos               Pos     `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Image             string  `json:"image"`
	Reference         string  `gorm:"not null" json:"reference"`
	Name              string  `gorm:"not null" json:"name"`
	Description       string  `gorm:"not null" json:"description"`
	UniteVente        string  `json:"unite_vente"`
	PrixVente         float64 `gorm:"not null" json:"prix_vente"`
	Tva               float64 `gorm:"default:0" json:"tva"`
	PrixAchat         float64 `gorm:"default:0" json:"prix_achat"`
	Remise            float64 `gorm:"default:0" json:"remise"`              // remise en pourcentage
	RemiseMinQuantite float64 `gorm:"default:0" json:"remise_min_quantite"` // remise en pourcentage pour la quantite minimale

	Stock          float64 `gorm:"default:0" json:"stock"`           // stock disponible
	StockEndommage float64 `gorm:"default:0" json:"stock_endommage"` // stock endommage
	Restitution    float64 `gorm:"default:0" json:"restitution"`     // stock restitution

	EntrepriseUUID string `json:"entreprise_uuid"`
	Signature      string `json:"signature"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	Stocks          []Stock          `gorm:"foreignKey:ProductUUID;references:UUID"` // Liste des stocks du produit
	StockEndommages []StockEndommage `gorm:"foreignKey:ProductUUID;references:UUID"` // Liste des stocks du produit
	CommadeLines    []CommandeLine   `gorm:"foreignKey:ProductUUID;references:UUID"` // Liste des stocks du produit
	Restitutions    []Restitution    `gorm:"foreignKey:ProductUUID;references:UUID"` // Liste des stocks du produit
}
