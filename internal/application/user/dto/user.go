package dto

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Anonymous bool   `json:"anonymous"`
}
