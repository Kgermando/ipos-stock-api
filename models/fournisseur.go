package models

import "gorm.io/gorm"

type Fournisseur struct {
	gorm.Model

	UUID           string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text
	Fullname       string `gorm:"not null" json:"fullname"`
	Telephone      string `gorm:"not null" json:"telephone"`
	Telephone2     string `json:"telephone2"`
	Email          string `json:"email"`
	Adress         string `json:"adress"`
	Entreprise     string `json:"entreprise"`
	WebSite        string `json:"website"`
	TypeFourniture string `json:"type_fourniture"`

	Signature      string `json:"signature"`
	CodeEntreprise uint64 `json:"code_entreprise"`

	Stocks []Stock `gorm:"foreignKey:FournisseurUUID"`
}
