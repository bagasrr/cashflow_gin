package response

type CategoryResponse struct {
	UserID  string `json:"user_id,omitempty"`
	GroupID string `json:"group_id,omitempty"`
	Name    string `json:"name" example:"Makanan"`
	Type    string `json:"type" example:"EXPENSE"`
}
