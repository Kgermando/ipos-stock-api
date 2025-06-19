package models

import (
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	UUID            string         `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	PosUUID         string         `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos             Pos            `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente
	Reference       uint64         `gorm:"not null" json:"reference"`          // Numero de reference du ravitaillement pour retrouver dans quel revitaillement le produit est endommagE
	ProductUUID     string         `gorm:"type:varchar(255);not null" json:"product_uuid"`
	Product         Product        `gorm:"foreignKey:ProductUUID;references:UUID"` // Produit associé
	Description     string         `json:"description"`
	Quantity        float64        `gorm:"not null" json:"quantity"`
	PrixAchat       float64        `gorm:"not null" json:"prix_achat"`
	DateExpiration  time.Time      `gorm:"not null" json:"date_expiration"`
	FournisseurUUID string         `gorm:"type:varchar(255);not null" json:"fournisseur_uuid"`
	Fournisseur     Fournisseur    `gorm:"foreignKey:FournisseurUUID;references:UUID"` // Fournisseur associé
	Signature       string         `json:"signature"`
	EntrepriseUUID  string         `json:"entreprise_uuid"`
	Sync            bool           `gorm:"default:false" json:"sync"`
}

type FournisseurStock struct {
	Name           string  `json:"name"`
	Telephone      string  `json:"telephone"`
	TypeFourniture string  `json:"type_fourniture"`
	TotalValue     float64 `json:"total_value"`
}
