package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Status bool `gorm:"not null" json:"status"`

	// Information de l'entreprise
	TypeEntreprise string `gorm:"not null" json:"type_entreprise"` // PME, GE, Particulier
	Name           string `gorm:"not null" json:"name"`
	Rccm           string `json:"rccm"`
	IdNat          string `json:"idnat"`
	NImpot         string `json:"nimpot"`
	Email          string `json:"email"`                     // Email officiel
	Telephone      string `gorm:"not null" json:"telephone"` // Telephone officiel
	Manager        string `gorm:"not null" json:"manager"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Country        string `json:"country"`

	// Information de payment
	Planid   string  `gorm:"not null" json:"planid"`               // ID du plan d'abonnement
	PlanName string  `gorm:"not null" json:"plan_name"`            // Nom du plan d'abonnement
	Amount   float64 `gorm:"not null" json:"amount"`               // Montant de l'abonnement
	Currency string  `gorm:"not null;default:CDF" json:"currency"` // Devise de l'abonnement, default CDF
	Duration int     `gorm:"not null" json:"duration"`             // Durée de l'abonnement en mois

	// Information de validation
	ValidateBy     string    `gorm:"not null" json:"validate_by"` // UUID de l'utilisateur qui a validé l'abonnement
	ValidationDate time.Time `json:"validation_date"`             // Date de validation de l'abonnement
	Notes          string    `json:"notes"`                       // Notes sur l'abonnement
	RejectedReason string    `json:"rejected_reason"`             // Raison du rejet de l'abonnement

}
