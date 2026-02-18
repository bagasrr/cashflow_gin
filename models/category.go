package models

import "github.com/google/uuid"

type Category struct {
	Base
	UserID      uuid.UUID     `gorm:"type:uuid" json:"user_id,omitempty"`
	GroupID     *uuid.UUID    `gorm:"type:uuid" json:"group_id,omitempty"`
	Name        string        `gorm:"type:varchar(100);unique" json:"name"`
	Type        string        `gorm:"type:varchar(20)" json:"type"`
	Transaction []Transaction `gorm:"foreignKey:CategoryID" json:"transactions,omitempty"`
}
