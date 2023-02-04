package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

var (
	ErrNoPrivKeyFile  = errors.New("private key file unreadable")
	ErrInvalidPrivKey = errors.New("RSA private key invalid")
	ErrNoPubKeyFile   = errors.New("public key file unreadable")
	ErrInvalidPubKey  = errors.New("RSA public key invalid")
)

func GetKey(key string) ([]byte, error) {
	if _, err := os.Stat(key); err == nil {
		return os.ReadFile(key)
	}
	return []byte(key), nil
}

func GetRSAPrivateKey(key string) (*rsa.PrivateKey, error) {
	if key == "" {
		return rsa.GenerateKey(rand.Reader, 4096)
	}

	prv, err := GetKey(key)
	if err != nil {
		return nil, ErrNoPrivKeyFile
	}
	prvKey, err := jwt.ParseRSAPrivateKeyFromPEM(prv)
	if err != nil {
		return nil, ErrInvalidPrivKey
	}
	return prvKey, nil
}

func GetRSAPublicKey(key string) (*rsa.PublicKey, error) {
	pub, err := GetKey(key)
	if err != nil {
		return nil, ErrNoPubKeyFile
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, ErrInvalidPubKey
	}
	return pubKey, nil
}
