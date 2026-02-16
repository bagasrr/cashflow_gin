package request

type CreateGroupRequest struct {
	Name        string   `json:"name" binding:"required" example:"Kelompok Keluarga"`
	Description string   `json:"description" example:"Kelompok untuk berbagi pengeluaran keluarga"`
	MemberIDs   []string `json:"member_ids"` // List user ID lain yg mau diajak (opsional)
}
