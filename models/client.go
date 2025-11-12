package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

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

	Commandes  []Commande  `gorm:"foreignKey:ClientUUID;references:UUID"` // Liste des commandes du client
	Livraisons []Livraison `gorm:"foreignKey:ClientUUID;references:UUID"` // Liste des livraisons du client
}
