package jwt

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"time"
)

var _ jwt.Claims = (*UserClaims)(nil)

type UserClaims struct {
	jwt.RegisteredClaims
	Username string          `json:"username,omitempty"`
	Email    string          `json:"email,omitempty"`
	Roles    []string        `json:"roles,omitempty"`
	OTP      bool            `json:"otp"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func RegisteredClaims(iss string, sub uuid.UUID, aud jwt.ClaimStrings, ttl time.Duration) jwt.RegisteredClaims {
	now := jwt.TimeFunc()

	return jwt.RegisteredClaims{
		Issuer:    iss,
		Subject:   sub.String(),
		Audience:  aud,
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
}

func (c UserClaims) UserID() uuid.UUID {
	return uuid.MustParse(c.Subject)
}

func (c UserClaims) Valid() error {
	vErr := &jwt.ValidationError{}
	now := jwt.TimeFunc()

	if !c.VerifyExpiresAt(now, true) {
		delta := now.Sub(c.ExpiresAt.Time)
		vErr.Inner = multierr.Append(vErr.Inner, fmt.Errorf("%s by %s", jwt.ErrTokenExpired, delta))
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if !c.VerifyIssuedAt(now, true) {
		vErr.Inner = multierr.Append(vErr.Inner, jwt.ErrTokenUsedBeforeIssued)
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if !c.VerifyNotBefore(now, true) {
		vErr.Inner = multierr.Append(vErr.Inner, jwt.ErrTokenNotValidYet)
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}
