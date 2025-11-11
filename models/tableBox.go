package models

import (
	"time"

	"gorm.io/gorm"
)

type TableBox struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PosUUID   string         `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos       Pos            `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente

	Name      string `json:"name"`
	Catergory string `json:"category"`
	Statut    string `json:"statut"`

	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	Commandes []Commande `gorm:"foreignKey:TableBoxUUID;references:UUID"` // Liste des lignes de commande

}
