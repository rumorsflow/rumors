package sys

import (
	"github.com/gowool/wool"
	"golang.org/x/exp/slog"
	"net/http"
)

type SignInDTO struct {
	Username string `json:"username,omitempty" validate:"required,min=3,max=254"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=64"`
}

type OtpDTO struct {
	Password string `json:"password,omitempty" validate:"required,numeric,len=6"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required,uuid4"`
}

type AuthActions struct {
	AuthService AuthService
	logger      *slog.Logger
}

func NewAuthActions(authService AuthService, logger *slog.Logger) *AuthActions {
	return &AuthActions{AuthService: authService, logger: logger}
}

func (a *AuthActions) SignIn(c wool.Ctx) error {
	var dto SignInDTO
	if err := c.Bind(&dto); err != nil {
		return err
	}

	session, err := a.AuthService.SignInByCredentials(c.Req().Context(), dto.Username, dto.Password)
	if err != nil {
		a.logger.Warn("sign in error", "err", err)

		return wool.NewErrBadRequest(nil)
	}

	return c.JSON(http.StatusOK, session)
}

func (a *AuthActions) OTP(c wool.Ctx) error {
	var dto OtpDTO
	if err := c.Bind(&dto); err != nil {
		return err
	}

	claims := GetClaims(c)

	session, err := a.AuthService.SignInByOTP(c.Req().Context(), claims.Username, dto.Password)
	if err != nil {
		a.logger.Warn("two factor auth error", "err", err)

		return wool.NewErrBadRequest(nil)
	}

	return c.JSON(http.StatusOK, session)
}

func (a *AuthActions) Refresh(c wool.Ctx) error {
	var dto RefreshTokenDTO
	if err := c.Bind(&dto); err != nil {
		return err
	}

	session, err := a.AuthService.SignInByRefreshToken(c.Req().Context(), dto.RefreshToken)
	if err != nil {
		a.logger.Warn("refresh token error", "err", err)

		return wool.NewErrBadRequest(nil)
	}

	return c.JSON(http.StatusOK, session)
}
