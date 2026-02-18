package response

import (
	"github.com/google/uuid"
)

type UserResponse struct {
	ID       string           `json:"id" example:"123e4567-e89b-12d3-a456-426655440000"`
	Username string           `json:"username" example:"john_doe"`
	Email    string           `json:"email" example:"john.doe@example.com"`
	UserRole string           `json:"user_role" example:"USER"`
	Wallets  []WalletResponse `json:"wallets,omitempty"`
}

type WalletResponse struct {
	ID               uuid.UUID             `json:"id" example:"123e4567-e89b-12d3-a456-426655440000"`
	Name             string                `json:"name" example:"Tabungan"`
	Balance          float64               `json:"balance" example:"1000"`
	GroupID          *uuid.UUID            `json:"group_id,omitempty"`
	Transactions     []TransactionResponse `json:"transactions,omitempty"`
	TransactionCount int64                 `json:"transaction_count" example:"5"`
}
