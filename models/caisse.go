package models

import "gorm.io/gorm"

type Caisse struct {
	gorm.Model

	UUID string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text 

	Name           string       `gorm:"not null" json:"name"` // Nom de la caisse 
	Signature      string       `json:"signature"`            // Signature de la transaction
	CodeEntreprise uint64       `json:"code_entreprise"`      // ID de l'entreprise
	PosUUID        string       `json:"pos_uuid"`             // ID du point de vente
	Pos            Pos          `gorm:"foreignKey:PosUUID"`   // Point de vente 
	Sync            bool       `gorm:"default:false" json:"sync"`

	
	Caisseitems    []CaisseItem `gorm:"foreignKey:CaisseUUID"`
}
