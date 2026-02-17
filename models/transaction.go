package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Base
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	WalletID uuid.UUID `gorm:"type:uuid;not null" json:"wallet_id"`

	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`

	Title            string    `gorm:"type:varchar(255)" json:"title"`
	Amount           float64   `gorm:"type:decimal(16,2)" json:"amount"`
	Description      string    `gorm:"type:text" json:"description"`
	Date             time.Time `json:"date"`
	TransactionCount int64     `gorm:"-:migration;->" json:"transaction_count"`

	User     User     `gorm:"foreignKey:UserID"`
	Wallet   Wallet   `gorm:"foreignKey:WalletID"`
	Category Category `gorm:"foreignKey:CategoryID"`
}
