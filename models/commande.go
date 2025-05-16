package models

import (
	"gorm.io/gorm"
)

type Commande struct {
	gorm.Model

	UUID           string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID        string `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos            Pos    `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Ncommande      uint64 `gorm:"not null" json:"ncommande"`          // Number Random
	Status         string `json:"status"`                             // Ouverte et Fermée
	ClientUUID     string `gorm:"type:varchar(255);not null" json:"client_uuid"`
	Client         Client `gorm:"foreignKey:ClientUUID;references:UUID"` // Client
	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	CommandeLines []CommandeLine `gorm:"foreignKey:CommandeUUID;references:UUID"` // Liste des lignes de commande
}
