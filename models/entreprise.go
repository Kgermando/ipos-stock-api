package models

import (
	"time"
	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type Entreprise struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TypeEntreprise string `gorm:"not null" json:"type_entreprise"` // PME, GE, Particulier
	Name           string `gorm:"not null" json:"name"` 
	Rccm           string `json:"rccm"`
	IdNat          string `json:"idnat"`
	NImpot         string `json:"nimpot"`
	Adresse        string `json:"adresse"`
	Email          string `json:"email"`                     // Email officiel
	Telephone      string `gorm:"not null" json:"telephone"` // Telephone officiel
	Manager        string `gorm:"not null" json:"manager"`
	Status         bool   `gorm:"not null" json:"status"`
	Currency       string `gorm:"not null;default:CDF" json:"currency"` // Devise de l'entreprise, default CDF

	Step int `gorm:"default:0" json:"step"` // Etape de l'entreprise dans le processus d'inscription

	TypeAbonnement string `json:"type_abonnement"` // Pack starter, business, pro, entreprise

	Users      []User       `gorm:"foreignKey:EntrepriseUUID;references:UUID"` // Liste des utilisateurs de l'entreprise
	Pos        []Pos        `gorm:"foreignKey:EntrepriseUUID;references:UUID"` // Liste des utilisateurs de l'entreprise
	Abonnement []Abonnement `gorm:"foreignKey:EntrepriseUUID;references:UUID"` // Liste des utilisateurs de l'entreprise
}

type EntrepriseInfos struct {
	UUID           string `json:"uuid"`            // Explicitly set type:text
	TypeEntreprise string `json:"type_entreprise"` // PME, GE, Particulier
	Name           string `json:"name"` 
	Rccm           string `json:"rccm"`
	IdNat          string `json:"idnat"`
	NImpot         string `json:"nimpot"`
	Adresse        string `json:"adresse"`
	Email          string `json:"email"`     // Email officiel
	Telephone      string `json:"telephone"` // Telephone officiel
	Manager        string `json:"manager"`
	Status         bool   `json:"status"`
	Currency       string `json:"currency"` // Devise de l'entreprise

	Step int `json:"step"` // Etape de l'entreprise dans le processus d'inscription

	TypeAbonnement string `json:"type_abonnement"` // Pack starter, business, pro, entreprise

	TotalUser       int `json:"total_user"`
	TotalPos        int `json:"total_pos"`
	TotalAbonnement int `json:"total_abonnement"`
}
