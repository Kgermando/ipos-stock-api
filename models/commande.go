package models

import (
	"time"

	"gorm.io/gorm"
)

type Commande struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PosUUID   string         `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos       Pos            `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente

	Ncommande string  `gorm:"not null" json:"ncommande"` // Number Random
	Status    string  `json:"status"`                    // Ouverte et Fermée
	TotalHt   float64 `gorm:"not null" json:"total_ht"`  // Total amount excluding tax
	TotalTva  float64 `gorm:"not null" json:"total_tva"` // Total tax amount
	TotalTtc  float64 `gorm:"not null" json:"total_ttc"` // Total amount including tax
	// TotalRemise    float64 `gorm:"not null" json:"total_remise"`       // Total discount amount
	// TotalGlobal    float64 `gorm:"not null" json:"total_global"`       // Total amount including all lines and taxes with reduction
	ClientUUID     string `gorm:"type:varchar(255);not null" json:"client_uuid"`
	Client         Client `gorm:"foreignKey:ClientUUID;references:UUID"` // Client
	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	TableBoxUUID string   `gorm:"type:varchar(255)" json:"table_box_uuid"`
	TableBox     TableBox `gorm:"foreignKey:TableBoxUUID;references:UUID"`

	CommandeLines []CommandeLine `gorm:"foreignKey:CommandeUUID;references:UUID"` // Liste des lignes de commande
}
