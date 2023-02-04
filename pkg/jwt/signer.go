package jwt

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v4"
)

type Signer interface {
	Sign(claims jwt.Claims) (string, error)
}

type signer struct {
	key *rsa.PrivateKey
}

func NewSigner(key *rsa.PrivateKey) Signer {
	return &signer{key: key}
}

func (s *signer) Sign(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodPS256, claims).SignedString(s.key)
}
