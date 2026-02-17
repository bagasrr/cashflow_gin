package models

import (
	"github.com/google/uuid"
)

type MembersRole int8

const (
	GroupAdmin       MembersRole = 1
	GroupParticipant MembersRole = 2
	GroupGuest       MembersRole = 3
)

func (g *MembersRole) String() string {
	switch *g {
	case GroupAdmin:
		return "ADMIN"
	case GroupParticipant:
		return "MEMBER"
	case GroupGuest:
		return "GUEST"
	default:
		return "UNKNOWN"
	}
}

type Group struct {
	Base
	Name        string `gorm:"type:varchar(200); not null" json:"group_name"`
	Description string `gorm:"type:text" json:"group_description"`

	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"group_owner_id"`
	MemberCount int64     `gorm:"-:migration;->" json:"member_count"`
	// WalletID uuid.UUID `gorm:"type:uuid;not null" json:"group_wallet_id"`

	Members []GroupMember `gorm:"foreignKey:GroupID" json:"members"`
	Wallet  []Wallet      `gorm:"foreignKey:GroupID" json:"wallets"`
}

type GroupMember struct {
	Base
	MembersRole MembersRole `gorm:"type:smallint;not null;default:2" json:"member_role"`
	// GroupID     uuid.UUID   `gorm:"type:uuid;not null" json:"group_id"`
	// UserID      uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	GroupID uuid.UUID `gorm:"type:uuid;not null;index:idx_member_group,unique" json:"group_id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_member_group,unique" json:"user_id"`

	User  User  `gorm:"foreignKey:UserID"`
	Group Group `gorm:"foreignKey:GroupID"`
}
