package request

type CreateWalletRequest struct {
	Name     string `json:"name" binding:"required,max=100" example:"Tabungan"`
	Currency string `json:"currency" binding:"required,len=3"` // e.g., IDR, USD
}
