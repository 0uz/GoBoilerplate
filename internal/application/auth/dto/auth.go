package auth

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AnonymousUserLoginRequest struct {
	Email string `json:"email" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AnonymousUserResponse struct {
	Email string `json:"email"`
}
