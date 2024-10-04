package dto

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefToken string `json:"refToken"`
}

type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type TokensObj struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
