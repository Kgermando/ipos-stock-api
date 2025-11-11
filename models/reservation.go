package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	UUID      string `gorm:"type:varchar(255);primary_key" json:"uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	PosUUID string `gorm:"type:varchar(255);not null" json:"pos_uuid"`
	Pos     Pos    `gorm:"foreignKey:PosUUID;references:UUID"` // Point de vente

	Table           string  `json:"table"`
	TableUUID       string  `json:"table_uuid"`
	ClientName      string  `json:"client_name"`
	ClientUUID      *string `json:"client_uuid,omitempty"` // Optionnel
	ReservationDate string  `json:"reservation_date"`
	ReservationTime string  `json:"reservation_time"`
	NumberOfGuests  int     `json:"number_of_guests"`
	Notes           *string `json:"notes,omitempty"` // Optionnel
	Status          string  `json:"status"`          // 'active', 'completed', 'cancelled'

	EntrepriseUUID string `json:"entreprise_uuid"`
	Signature      string `json:"signature"`
	Sync           bool   `gorm:"default:false" json:"sync"`
}
