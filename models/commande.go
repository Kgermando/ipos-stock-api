package models

import (
	"gorm.io/gorm"
)

type Commande struct {
	gorm.Model

	UUID           string         `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	PosUUID        string         `json:"pos_uuid"`
	Pos            Pos            `gorm:"foreignKey:PosUUID"`
	Ncommande      uint64         `gorm:"not null" json:"ncommande"` // Number Random
	Status         string         `json:"status"`                    // Ouverte et Fermée
	ClientUUID     string         `json:"client_uuid"`
	Client         Client         `gorm:"foreignKey:ClientUUID"`
	Signature      string         `json:"signature"`
	CodeEntreprise uint64         `json:"code_entreprise"`
	Sync            bool       `gorm:"default:false" json:"sync"`

	
	CommandeLines  []CommandeLine `gorm:"foreignKey:CommandeUUID"`
}
