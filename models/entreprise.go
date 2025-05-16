package models

import (

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type Entreprise struct {
	gorm.Model

	UUID           string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	TypeEntreprise string `gorm:"not null" json:"type_entreprise"`       // PME, GE, Particulier
	Name           string `gorm:"not null" json:"name"`
	Code           uint64 `gorm:"not null" json:"code"` // Code entreprise
	Rccm           string `json:"rccm"`
	IdNat          string `json:"idnat"`
	NImpot         string `json:"nimpot"`
	Adresse        string `json:"adresse"`
	Email          string `json:"email"`                     // Email officiel
	Telephone      string `gorm:"not null" json:"telephone"` // Telephone officiel
	Manager        string `gorm:"not null" json:"manager"`
	Status         bool   `gorm:"not null" json:"status"`
	Currency       string `gorm:"not null;default:CDF" json:"currency"` // Devise de l'entreprise, default CDF
	TypeAbonnement string `json:"type_abonnement"`
	Signature      string `json:"signature"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	Users      []User       `gorm:"foreignKey:EntrepriseUUID;references:UUID"` // Liste des utilisateurs de l'entreprise
	Pos        []Pos        `gorm:"foreignKey:EntrepriseUUID;references:UUID"` // Liste des utilisateurs de l'entreprise
	Abonnement []Abonnement `gorm:"foreignKey:EntrepriseUUID;references:UUID"` // Liste des utilisateurs de l'entreprise
}

type EntrepriseInfos struct {
	ID              uint   `json:"id"`
	UUID            string `json:"uuid"`            // Explicitly set type:text
	TypeEntreprise  string `json:"type_entreprise"` // PME, GE, Particulier
	Name            string `json:"name"`
	Code            uint64 `json:"code"` // Code entreprise
	Rccm            string `json:"rccm"`
	IdNat           string `json:"idnat"`
	NImpot          string `json:"nimpot"`
	Adresse         string `json:"adresse"`
	Email           string `json:"email"`     // Email officiel
	Telephone       string `json:"telephone"` // Telephone officiel
	Manager         string `json:"manager"`
	Status          bool   `json:"status"`
	Currency        string `json:"currency"` // Devise de l'entreprise
	TypeAbonnement  string `json:"type_abonnement"`
	Signature       string `json:"signature"`
	TotalUser       int    `json:"total_user"`
	TotalPos        int    `json:"total_pos"`
	TotalAbonnement int    `json:"total_abonnement"`
}
