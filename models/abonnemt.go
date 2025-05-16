package models

import (
	"gorm.io/gorm"
)

type Abonnement struct {
	gorm.Model

	UUID           string     `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	EntrepriseUUID string     `gorm:"type:varchar(255);not null" json:"entreprise_uuid"`
	Entreprise     Entreprise `gorm:"foreignKey:EntrepriseUUID;references:UUID"`
	Bouquet        string     `gorm:"not null" json:"bouquet"` // Premium, Platinuim, Fremuim
	Montant        float64    `gorm:"not null" json:"montant"`
	MoyenPayment   string     `gorm:"not null" json:"moyen_payment"`
	Signature      string     `json:"signature"`
	Sync           bool       `gorm:"default:false" json:"sync"`
}
