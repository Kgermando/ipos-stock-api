package models

import (
	"gorm.io/gorm"
)

type Pos struct {
	gorm.Model

	UUID           string     `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	EntrepriseUUID string     `json:"entreprise_uuid"`
	Entreprise     Entreprise `gorm:"foreignKey:EntrepriseUUID"`
	Name           string     `gorm:"not null" json:"name"`
	Adresse        string     `json:"adresse"`
	Email          string     `json:"email"`
	Telephone      string     `json:"telephone"`
	Manager        string     `gorm:"not null" json:"manager"`
	Status         bool       `gorm:"not null" json:"status"` // Actif ou Inactif
	CodeEntreprise uint64     `json:"code_entreprise"`
	Signature      string     `json:"signature"`

	Stocks    []Stock    `gorm:"foreignKey:PosID" json:"stocks"`
	Commandes []Commande `gorm:"foreignKey:PosID" json:"commandes"`
}
