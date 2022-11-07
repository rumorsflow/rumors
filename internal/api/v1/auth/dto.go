package auth

import "github.com/cristalhq/jwt/v4"

type SessionRequest struct {
	Username string `json:"username,omitempty" validate:"required,max=254"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=64"`
}

type OtpRequest struct {
	Password string `json:"password,omitempty" validate:"required,numeric,len=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,uuid4"`
}

type SessionResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username  string            `json:"username,omitempty"`
	Email     string            `json:"email,omitempty"`
	Roles     []string          `json:"roles,omitempty"`
	Two2FA    bool              `json:"two_fa,omitempty"`
	Providers map[string]string `json:"providers,omitempty"`
}
