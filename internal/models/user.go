package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Provider string

const (
	Google   Provider = "GOOGLE"
	Facebook Provider = "FACEBOOK"
	Linkedin Provider = "LINKEDIN"
	GitHub   Provider = "GITHUB"
)

type ProviderData struct {
	Name   Provider `json:"name,omitempty" bson:"name,omitempty"`
	Value  string   `json:"value,omitempty" bson:"value,omitempty"`
	Delete bool     `json:"-" bson:"-"`
}

func (p *ProviderData) Is(provider Provider) bool {
	return p.Name == provider
}

type User struct {
	Id          string         `json:"id,omitempty" bson:"_id,omitempty"`
	Username    string         `json:"username,omitempty" bson:"username,omitempty"`
	Email       string         `json:"email,omitempty" bson:"email,omitempty"`
	Password    string         `json:"-" bson:"password,omitempty"`
	OTPSecret   []byte         `json:"-" bson:"otp_secret,omitempty"`
	Roles       []string       `json:"roles,omitempty" bson:"roles,omitempty"`
	DeleteRoles []string       `json:"-" bson:"-"`
	Providers   []ProviderData `json:"providers,omitempty" bson:"providers,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (u *User) HasPassword() bool {
	return u.Password != ""
}

func (u *User) Secret() string {
	return string(u.OTPSecret)
}

func (u *User) HasProvider(provider Provider) bool {
	return lo.ContainsBy(u.Providers, func(p ProviderData) bool {
		return p.Is(provider)
	})
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}

func (u *User) GenerateOTPSecret(size uint) (err error) {
	u.OTPSecret, err = generateSecret(size)
	return err
}

func (u *User) GeneratePasswordHash() error {
	pwd, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password, err: %w", err)
	}
	return string(hash), nil
}

func generateSecret(size uint) ([]byte, error) {
	s := make([]byte, size)
	_, err := rand.Reader.Read(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
