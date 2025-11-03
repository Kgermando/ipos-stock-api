package models

import (
	"time"

	"gorm.io/gorm"
)

type Plat struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PosUUID   string         `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos       Pos            `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente

	Image       string  `json:"image"`
	Reference   string  `gorm:"not null" json:"reference"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	Categorie   string  `json:"categorie"`
	Prix        float64 `gorm:"not null" json:"prix"`
	Tva         float64 `gorm:"default:0" json:"tva"`
	Remise      float64 `gorm:"default:0" json:"remise"` // remise en pourcentage

	// Spécifique aux plats - pas de gestion de stock quantifiable
	IsAvailable bool `gorm:"default:true" json:"is_available"` // Disponibilité du plat

	EntrepriseUUID string `json:"entreprise_uuid"`
	Signature      string `json:"signature"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	// Relations avec les commandes
	CommadeLines []CommandeLine `gorm:"foreignKey:PlatUUID;references:UUID"` // Liste des commandes contenant ce plat
}
