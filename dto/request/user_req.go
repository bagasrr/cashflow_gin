package request

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,max=100" example:"john_doe"`
	Email    string `json:"email" binding:"required,max=100" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,max=255" example:"password123"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required" example:"johndoeganteng"`
}
