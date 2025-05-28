package models

import (
	"gorm.io/gorm"
)

type Commande struct {
	gorm.Model

	UUID           string  `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID        string  `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos            Pos     `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Ncommande      uint64  `gorm:"not null" json:"ncommande"`          // Number Random
	Status         string  `json:"status"`                             // Ouverte et Fermée
	TotalHt        float64 `gorm:"not null" json:"total_ht"`           // Total amount excluding tax
	TotalTva       float64 `gorm:"not null" json:"total_tva"`          // Total tax amount
	TotalTtc       float64 `gorm:"not null" json:"total_ttc"`          // Total amount including tax
	TotalRemise    float64 `gorm:"not null" json:"total_remise"`       // Total discount amount
	TotalGlobal    float64 `gorm:"not null" json:"total_global"`       // Total amount including all lines and taxes with reduction
	ClientUUID     string  `gorm:"type:varchar(255);not null" json:"client_uuid"`
	Client         Client  `gorm:"foreignKey:ClientUUID;references:UUID"` // Client
	Signature      string  `json:"signature"`
	EntrepriseUUID string  `json:"entreprise_uuid"`
	Sync           bool    `gorm:"default:false" json:"sync"`

	CommandeLines []CommandeLine `gorm:"foreignKey:CommandeUUID;references:UUID"` // Liste des lignes de commande
}
