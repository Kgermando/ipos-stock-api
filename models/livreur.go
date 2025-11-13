package models

import (
	"time"

	"gorm.io/gorm"
)

type Livreur struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	TypeLivreur string `json:"type_livreur"` // 'Entreprise' ou 'Organisation' ou 'Particulier'
	Name        string `gorm:"not null" json:"name"`
	Telephone   string `gorm:"not null" json:"telephone"`
	Email       string `json:"email"`
	Adresse     string `json:"adresse"`
	Manager     string `gorm:"not null" json:"manager"`

	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`

	PosUUID string `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos     Pos    `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Sync    bool   `gorm:"default:false" json:"sync"`

	Livraisons []Livraison `gorm:"foreignKey:LivreurUUID;references:UUID"` // Liste des livraisons du livreur
}
