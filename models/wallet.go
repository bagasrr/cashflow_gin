package models

import "github.com/google/uuid"

type Wallet struct {
	Base
	UserID  *uuid.UUID `gorm:"type:uuid;" json:"user_id"`
	GroupID *uuid.UUID `gorm:"type:uuid" json:"group_id,omitempty"`
	Name    string     `gorm:"type:varchar(100)" json:"name"`
	// UBAH MUTLAK: Pakai int64 dan bigint.
	Balance      int64         `gorm:"type:bigint;default:0" json:"balance"`
	Currency     string        `gorm:"type:varchar(10);default:'IDR'" json:"currency"`
	Transactions []Transaction `gorm:"foreignKey:WalletID" json:"transactions,omitempty"`
	Groups       *Group        `gorm:"foreignKey:GroupID" json:"groups,omitempty"`

	TransactionCount int64 `gorm:"-:migration;->" json:"transaction_count"`
}
