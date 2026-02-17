package response

import "time"

type TxWithWallet struct {
	TransactionID string    `json:"transaction_id"`
	Title         string    `json:"title" example:"Gaji Bulanan"`
	Amount        float64   `json:"amount" example:"500"`
	Description   string    `json:"description" example:"Gaji bulan Januari 2026"`
	Date          time.Time `json:"date" example:"2026-01-31T00:00:00Z" format:"date-time"`

	User     UserResponse     `json:"user"`
	Category CategoryResponse `json:"category"`
	Wallet   WalletResponse   `json:"wallet"`
}

type TransactionResponse struct {
	ID          string           `json:"id" example:"123e4567-e89b-12d3-a456-426655440000"`
	Title       string           `json:"title" example:"Gaji Bulanan"`
	Amount      float64          `json:"amount" example:"500"`
	Description string           `json:"description" example:"Gaji bulan Januari 2026"`
	Date        time.Time        `json:"date" example:"2026-01-31T00:00:00Z" format:"date-time"`
	Category    CategoryResponse `json:"category"`
	User        UserResponse     `json:"user"`
}
