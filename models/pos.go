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
	Signature      string     `json:"signature"`
	CodeEntreprise uint64     `json:"code_entreprise"`

	Users        []User        `gorm:"foreignKey:PosUUID"`
	Products     []Product     `gorm:"foreignKey:PosUUID"`
	Commandes    []Commande    `gorm:"foreignKey:PosUUID"`
	Clients      []Client      `gorm:"foreignKey:PosUUID"`
	Fournisseurs []Fournisseur `gorm:"foreignKey:PosUUID"`
}
