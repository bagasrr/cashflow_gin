package models

import (
	"strings"
	"time"
)

type UserRole int8

const (
	RoleAdmin     UserRole = 1
	RoleModerator UserRole = 2
	RoleUser      UserRole = 3
)

func (u *UserRole) String() string {
	switch *u {
	case RoleAdmin:
		return "admin"
	case RoleModerator:
		return "moderator"
	case RoleUser:
		return "user"
	default:
		return "unknown"
	}
}
func ParseUserRole(role string) UserRole {
	switch strings.ToLower(role) {
	case "admin":
		return RoleAdmin
	case "moderator":
		return RoleModerator
	case "user":
		return RoleUser
	default:
		return RoleUser // Fallback aman, default selalu jadi user biasa
	}
}

type User struct {
	Base
	Username              string     `gorm:"type:varchar(100);unique" json:"username"`
	Email                 string     `gorm:"type:varchar(100);unique" json:"email"`
	Password              string     `gorm:"type:varchar(255)" json:"-"` // Hide password dari JSON
	UserRole              UserRole   `gorm:"type:smallint" json:"user_role" default:"3"`
	SubscriptionPlan      string     `gorm:"type:varchar(100)" json:"subscription_plan"`
	SubscriptionExpiredAt *time.Time `json:"subscription_expired_at" `

	// Relations
	Wallets []Wallet `gorm:"foreignKey:UserID"`
}
