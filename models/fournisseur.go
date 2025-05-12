package models

import "gorm.io/gorm"

type Fournisseur struct {
	gorm.Model

	UUID           string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	EntrepriseName string `gorm:"not null" json:"entreprise_name"`
	Rccm           string `json:"rccm"`
	IdNat          string `json:"idnat"`
	NImpot         string `json:"nimpot"`
	Adresse        string `json:"adresse"`
	Email          string `json:"email"`                     // Email officiel
	Telephone      string `gorm:"not null" json:"telephone"` // Telephone officiel
	Manager        string `gorm:"not null" json:"manager"`
	WebSite        string `json:"website"`
	TypeFourniture string `json:"type_fourniture"`

	Signature      string `json:"signature"`
	CodeEntreprise uint64 `json:"code_entreprise"`

	PosUUID string `json:"pos_uuid"`
	Pos     Pos    `gorm:"foreignKey:PosUUID"`
	Sync            bool       `gorm:"default:false" json:"sync"`

	Stocks []Stock `gorm:"foreignKey:FournisseurUUID"`
}
