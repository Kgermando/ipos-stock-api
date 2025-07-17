package models

import (
	"time"

	"gorm.io/gorm"
)

type Abonnement struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	EntrepriseUUID string     `gorm:"type:varchar(255);not null" json:"entreprise_uuid"`
	Entreprise     Entreprise `gorm:"foreignKey:EntrepriseUUID;references:UUID"`
	Montant        float64    `gorm:"not null" json:"montant"`
	MoyenPayment   string     `gorm:"not null" json:"moyen_payment"`            // Ex: "card", "bank_transfer" mpesa, etc.
	Duree          int        `gorm:"not null" json:"duree"`                    // Duration in months ex. 14 months
	Statut         string     `gorm:"not null;default:'pending'" json:"statut"` // pending, active, suspended, cancelled
	Signature      string     `json:"signature"`
}
