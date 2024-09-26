package presenter

type TokenResponse struct {
	AccessToken  string `json:"token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
