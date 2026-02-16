package models

import "github.com/google/uuid"

type Wallet struct {
	Base
	UserID      *uuid.UUID    `gorm:"type:uuid;" json:"user_id"`
	GroupID     *uuid.UUID    `gorm:"type:uuid" json:"group_id,omitempty"` // Nullable untuk wallet pribadi
	Name        string        `gorm:"type:varchar(100)" json:"name"`
	Balance     float64       `gorm:"type:decimal(16,2);default:0" json:"balance"`
	Currency    string        `gorm:"type:varchar(10);default:'IDR'" json:"currency"`
	Transaction []Transaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
	Groups      *Group        `gorm:"foreignKey:GroupID" json:"groups,omitempty"`
}
