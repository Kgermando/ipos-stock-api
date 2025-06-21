package models

import (
	"time"

	"gorm.io/gorm"
)

type Abonnement struct {
	UUID           string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	EntrepriseUUID string         `gorm:"type:varchar(255);not null" json:"entreprise_uuid"`
	Entreprise     Entreprise     `gorm:"foreignKey:EntrepriseUUID;references:UUID"`
	Bouquet        string         `gorm:"not null" json:"bouquet"` // Premium, Platinuim, Fremuim
	Montant        float64        `gorm:"not null" json:"montant"`
	MoyenPayment   string         `gorm:"not null" json:"moyen_payment"`
	Signature      string         `json:"signature"`
	Sync           bool           `gorm:"default:false" json:"sync"`
}
