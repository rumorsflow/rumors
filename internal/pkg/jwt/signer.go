package jwt

import (
	"crypto/rsa"
	"github.com/cristalhq/jwt/v4"
)

type Signer interface {
	Sign(claims any) (string, error)
}

type signer struct {
	builder *jwt.Builder
}

func NewSigner(privKey *rsa.PrivateKey) (Signer, error) {
	s, err := jwt.NewSignerRS(jwt.RS256, privKey)
	if err != nil {
		return nil, err
	}
	return &signer{builder: jwt.NewBuilder(s)}, nil
}

func (s *signer) Sign(claims any) (string, error) {
	token, err := s.builder.Build(claims)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}
