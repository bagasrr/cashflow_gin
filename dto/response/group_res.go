package response

// dto/request/group_request.go

// dto/response/group_response.go
type GroupResponse struct {
	ID           string                `json:"id" example:"123e4567-e89b-12d3-a456-426655440000"`
	Name         string                `json:"name" example:"Kelompok Keluarga"`
	Description  string                `json:"description" example:"Kelompok untuk berbagi pengeluaran keluarga"`
	Wallet       WalletResponse        `json:"wallet"` // Group pasti punya wallet
	Members      []GroupMemberResponse `json:"members,omitempty"`
	TotalMembers int64                 `json:"total_members,omitempty" example:"5"`
}

type GroupMemberResponse struct {
	ID       string `json:"id" example:"123e4567-e89b-12d3-a456-426655440000"`
	UserID   string `json:"user_id"` // ID User aslinya
	Username string `json:"username" example:"john_doe"`
	Role     string `json:"role" example:"ADMIN"`
}
