package models

import (
	"time"

	"gorm.io/gorm"
)

type Livraison struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	PosUUID string `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos     Pos    `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente

	ClientUUID string `gorm:"type:varchar(255);not null" json:"client_uuid"`
	Client     Client `gorm:"foreignKey:ClientUUID;references:UUID"` // Client

	LivreurUUID string      `gorm:"type:varchar(255);not null" json:"livreur_uuid"`
	Livreur     Livreur `gorm:"foreignKey:LivreurUUID;references:UUID"` // Livreur

	ZoneUUID string `gorm:"type:varchar(255);not null" json:"zone_uuid"`
	Zone     Zone   `gorm:"foreignKey:ZoneUUID;references:UUID"` // Zone de livraison

	Statut string `gorm:"not null" json:"statut"` // En cours, Effectuée, Annulée

	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`
	Sync           bool   `gorm:"default:false" json:"sync"`

	Commandes []Commande `gorm:"foreignKey:LivraisonUUID;references:UUID"` // Liste des commandes associées à la livraison
}
