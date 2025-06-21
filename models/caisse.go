package models

import (
	"time"

	"gorm.io/gorm"
)

type Caisse struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name           string  `gorm:"not null" json:"name"`                       // Nom de la caisse
	MontantEntre   float64 `gorm:"default:0" json:"montant_entre"`             // Montant d'entrée
	MontantSorti   float64 `gorm:"default:0" json:"montant_sorti"`             // Montant de sortie
	MontantDebut   float64 `gorm:"default:0" json:"montant_debut"`             // Montant de début
	Signature      string  `json:"signature"`                                  // Signature de la transaction
	EntrepriseUUID string  `json:"entreprise_uuid"`                            // ID de l'entreprise
	PosUUID        string  `gorm:"type:varchar(255);not null" json:"pos_uuid"` // ID du point de vente
	Pos            Pos     `gorm:"foreignKey:PosUUID;references:UUID"`         // Point de vente
	Sync           bool    `gorm:"default:false" json:"sync"`

	Caisseitems []CaisseItem `gorm:"foreignKey:CaisseUUID;references:UUID"` // Liste des items de la caisse
}
