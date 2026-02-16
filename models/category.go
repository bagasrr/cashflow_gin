package models

import "github.com/google/uuid"

type Category struct {
	Base
	UserID      uuid.UUID     `gorm:"type:uuid" json:"user_id"` // Nullable (Global category) atau Specific User
	Name        string        `gorm:"type:varchar(100)" json:"name"`
	Type        string        `gorm:"type:varchar(20)" json:"type"` // INCOME / EXPENSE
	Transaction []Transaction `gorm:"foreignKey:CategoryID" json:"transactions,omitempty"`
}
