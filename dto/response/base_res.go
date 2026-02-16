package response

type BaseResponse struct {
	Status  bool        `json:"status" example:"true"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty"` // omitempty: sembunyikan jika nil
	Errors  interface{} `json:"errors,omitempty"`
}
