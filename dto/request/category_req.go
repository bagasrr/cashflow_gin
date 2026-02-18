package request

type CreateCategoryRequest struct {
	Name    string `json:"name" binding:"required,max=100" example:"Makanan"`
	Type    string `json:"type" binding:"required,oneof=INCOME EXPENSE" example:"EXPENSE"`
	UserID  string `json:"user_id" example:"123e4567-e89b-12d3-a456-426655440000"`
	GroupID string `json:"group_id,omitempty" example:"123e4567-e89b-12d3-a456-426655440000"`
}
