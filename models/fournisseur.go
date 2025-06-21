package models

import (
	"time"

	"gorm.io/gorm"
)

type Fournisseur struct {
	UUID           string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	EntrepriseName string         `gorm:"not null" json:"entreprise_name"`
	Rccm           string         `json:"rccm"`
	IdNat          string         `json:"idnat"`
	NImpot         string         `json:"nimpot"`
	Adresse        string         `json:"adresse"`
	Email          string         `json:"email"`                     // Email officiel
	Telephone      string         `gorm:"not null" json:"telephone"` // Telephone officiel
	Manager        string         `gorm:"not null" json:"manager"`
	WebSite        string         `json:"website"`
	TypeFourniture string         `json:"type_fourniture"`

	Signature      string `json:"signature"`
	EntrepriseUUID string `json:"entreprise_uuid"`

	PosUUID string `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos     Pos    `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Sync    bool   `gorm:"default:false" json:"sync"`

	Stocks []Stock `gorm:"foreignKey:FournisseurUUID;references:UUID"` // Liste des stocks du fournisseur
}
