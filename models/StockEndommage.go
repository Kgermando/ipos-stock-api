package models

import (
	"time"

	"gorm.io/gorm"
)

type StockEndommage struct {
	gorm.Model

	UUID           string    `gorm:"type:text;not null;unique" json:"uuid"`
	PosUUID        string    `json:"pos_uuid"`
	Pos            Pos       `gorm:"foreignKey:PosUUID"`
	StockUUID      string    `gorm:"not null" json:"stock_uuid"`
	Stock          Stock     `gorm:"foreignKey:StockUUID"`
	Raison         string    `json:"raison"` // Raison de l'endommagement
	DateEndommage  time.Time `json:"date_endommage"`
	Quantity        float64      `gorm:"not null" json:"quantity"`
	CodeEntreprise uint64    `json:"code_entreprise"`
}
