package models

import (
	"time"

	"gorm.io/gorm"
)

type Pos struct {
	UUID           string         `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	EntrepriseUUID string         `gorm:"type:varchar(255);not null" json:"entreprise_uuid"`
	Entreprise     Entreprise     `gorm:"foreignKey:EntrepriseUUID;references:UUID"`
	Name           string         `gorm:"not null" json:"name"`
	Adresse        string         `json:"adresse"`
	Email          string         `json:"email"`
	Telephone      string         `json:"telephone"`
	Manager        string         `gorm:"not null" json:"manager"`
	Status         bool           `gorm:"not null" json:"status"` // Actif ou Inactif
	Signature      string         `json:"signature"`
	CodeEntreprise uint64         `json:"code_entreprise"`
	Sync           bool           `gorm:"default:false" json:"sync"`

	Users        []User        `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Products     []Product     `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Commandes    []Commande    `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Clients      []Client      `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Fournisseurs []Fournisseur `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
}
