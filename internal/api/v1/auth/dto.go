package auth

type SignUpRequest struct {
	Username       string `json:"username,omitempty" validate:"required,min=3,max=254"`
	Email          string `json:"email,omitempty" validate:"required,email,min=3,max=254"`
	Password       string `json:"password" validate:"required,eqfield=RepeatPassword,min=8,max=64"`
	RepeatPassword string `json:"repeat_password" validate:"required,min=8,max=64"`
}

type SignUpResponse struct {
	Uri string `json:"uri"`
	Qr  string `json:"qr"`
}

type SignInRequest struct {
	Username string `json:"username,omitempty" validate:"required,min=3,max=254"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=64"`
}

type OtpRequest struct {
	Password string `json:"password,omitempty" validate:"required,numeric,len=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,uuid4"`
}
