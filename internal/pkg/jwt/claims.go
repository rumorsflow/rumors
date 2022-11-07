package jwt

import (
	"errors"
	"github.com/cristalhq/jwt/v4"
	"github.com/rumorsflow/rumors/internal/models"
	"time"
)

var (
	ErrTokenInvalidAudience  = errors.New("token has invalid audience")
	ErrTokenExpired          = errors.New("token is expired")
	ErrTokenUsedBeforeIssued = errors.New("token used before issued")
	ErrTokenNotValidYet      = errors.New("token is not valid yet")
)

const Audience = "Rumors"

type UserClaims struct {
	jwt.RegisteredClaims
	Username  string                `json:"username,omitempty"`
	Email     string                `json:"email,omitempty"`
	Roles     []models.Role         `json:"roles,omitempty"`
	Providers []models.ProviderData `json:"providers,omitempty"`
	Metadata  any                   `json:"metadata,omitempty"`
}

func (uc UserClaims) Validate() error {
	if valid := uc.IsForAudience(Audience); !valid {
		return ErrTokenInvalidAudience
	}
	if valid := uc.IsValidAt(time.Now()); !valid {
		return ErrTokenExpired
	}
	if valid := uc.IsValidIssuedAt(time.Now()); !valid {
		return ErrTokenUsedBeforeIssued
	}
	if valid := uc.IsValidNotBefore(time.Now()); !valid {
		return ErrTokenNotValidYet
	}
	return nil
}

func RegisteredClaims(id string, ttl time.Duration) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ID:        id,
		Audience:  jwt.Audience{Audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
}
