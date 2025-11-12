package models

import (
	"time"

	"gorm.io/gorm"
)

type Zone struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	PosUUID string `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos     Pos    `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente

	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	Livraisons []Livraison `gorm:"foreignKey:ZoneUUID;references:UUID"` // Liste des livraisons dans cette zone
}
