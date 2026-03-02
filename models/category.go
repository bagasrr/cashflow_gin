package models

import "github.com/google/uuid"

type Category struct {
	Base
	// 1. Tambahkan index gabungan bernama 'idx_user_category_name'
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_category_name" json:"user_id,omitempty"`

	GroupID *uuid.UUID `gorm:"type:uuid" json:"group_id,omitempty"`

	// 2. Masukkan Name ke dalam index gabungan yang sama, HAPUS 'unique' bawaannya
	Name string `gorm:"type:varchar(100);uniqueIndex:idx_user_category_name" json:"name"`

	Type        string        `gorm:"type:varchar(20)" json:"type"`
	Transaction []Transaction `gorm:"foreignKey:CategoryID" json:"transactions,omitempty"`
}
