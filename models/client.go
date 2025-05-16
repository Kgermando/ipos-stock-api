package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model

	UUID string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text

	Fullname   string `gorm:"not null" json:"fullname"`
	Telephone  string `gorm:"not null" json:"telephone"`
	Telephone2 string `json:"telephone2"`
	Email      string `json:"email"`
	Adress     string `json:"adress"`
	// Birthday     string `json:"birthday"`
	Organisation string `json:"organisation"`
	WebSite      string `json:"website"`

	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`
	PosUUID        string `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos            Pos    `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Sync           bool   `gorm:"default:false" json:"sync"`

	Commandes []Commande `gorm:"foreignKey:ClientUUID;references:UUID"` // Liste des commandes du client
}
