package models

import "gorm.io/gorm"

type CaisseItem struct {
	gorm.Model

	UUID       string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	CaisseUUID string `gorm:"type:varchar(255);not null" json:"caisse_uuid"`
	Caisse     Caisse `gorm:"foreignKey:CaisseUUID;references:UUID"` // Caisse associée

	TypeTransaction string  `gorm:"not null" json:"type_transaction"` // Entrée ou Sortie
	Montant         float64 `gorm:"not null" json:"montant"`          // Montant de la transaction
	Libelle         string  `json:"libelle"`                          // Description de la transaction
	Reference       string  `json:"reference"`                        // Nombre aleatoire
	Signature       string  `json:"signature"`                        // Signature de la transaction
	EntrepriseUUID  string  `json:"entreprise_uuid"`
	PosUUID         string  `gorm:"type:varchar(255);not null" json:"pos_uuid"` // ID du point de vente
	Pos             Pos     `gorm:"foreignKey:PosUUID;references:UUID"`         // Point de vente
	Sync            bool    `gorm:"default:false" json:"sync"`                  // ID de l'entreprise
}
