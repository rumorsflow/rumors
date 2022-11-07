package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/cristalhq/jwt/v4"
)

var (
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")
	ErrTokenInvalidClaims      = errors.New("token has invalid claims")
)

type Verifier interface {
	Parse(token string) (*jwt.Token, UserClaims, error)
}

type verifier struct {
	jwt.Verifier
}

func NewVerifier(pubKey *rsa.PublicKey) (Verifier, error) {
	v, err := jwt.NewVerifierRS(jwt.RS256, pubKey)
	if err != nil {
		return nil, err
	}
	return &verifier{Verifier: v}, nil
}

func (v *verifier) Parse(token string) (*jwt.Token, UserClaims, error) {
	var claims UserClaims
	t, err := jwt.Parse([]byte(token), v)
	if err != nil {
		return nil, claims, ErrInvalidSigningAlgorithm
	}
	claims, err = Unmarshal(t.Claims())
	if err == nil {
		err = claims.Validate()
	}
	return t, claims, err
}

func Unmarshal(claims json.RawMessage) (UserClaims, error) {
	var uc UserClaims
	if err := json.Unmarshal(claims, &uc); err != nil {
		return uc, ErrTokenInvalidClaims
	}
	return uc, nil
}
