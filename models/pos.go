package models

import (
	"time"

	"gorm.io/gorm"
)

type Pos struct {
	UUID           string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
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

	Users           []User          `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Products        []Product       `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Caisses         []Caisse        `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des caisses associées au point de vente
	CaisseItems     []CaisseItem    `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des items de caisse associés au point de vente
	Commandes       []Commande      `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Clients         []Client        `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	Fournisseurs    []Fournisseur   `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des utilisateurs du point de vente
	CommandeLines   []CommandeLine  `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des lignes de commande associées au point de vente
	Stocks          []Stock         `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des stocks associés au point de vente
	StockEndommages []StockEndommage `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des stocks endommagés associés au point de vente
	Restitutions    []Restitution   `gorm:"foreignKey:PosUUID;references:UUID"` // Liste des restitutions associées au point de vente
}
