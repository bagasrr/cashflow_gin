package models

import "github.com/google/uuid"

type Category struct {
	Base
	// 1. Tambahkan index gabungan bernama 'idx_user_category_name'
	GroupID *uuid.UUID `gorm:"type:uuid" json:"group_id,omitempty"`
	// UserID  uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_user_category_name" json:"user_id,omitempty"`
	// Name        string        `gorm:"type:varchar(100);uniqueIndex:idx_user_category_name" json:"name"`
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_category_name,where:deleted_at IS NULL" json:"user_id,omitempty"`
	Name   string    `gorm:"type:varchar(100);uniqueIndex:idx_user_category_name,where:deleted_at IS NULL" json:"name"`
	// 2. Masukkan Name ke dalam index gabungan yang sama, HAPUS 'unique' bawaannya
	Type        string        `gorm:"type:varchar(20)" json:"type"`
	Transaction []Transaction `gorm:"foreignKey:CategoryID" json:"transactions,omitempty"`
}
