package models

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/pquerna/otp/totp"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type (
	Role     string
	Provider string
)

const (
	Google   Provider = "GOOGLE"
	Facebook Provider = "FACEBOOK"
	Linkedin Provider = "LINKEDIN"
	GitHub   Provider = "GITHUB"

	AdminRole Role = "ROLE_ADMIN"
	UserRole  Role = "ROLE_USER"
)

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

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
	Roles       []Role         `json:"roles,omitempty" bson:"roles,omitempty"`
	DeleteRoles []Role         `json:"-" bson:"-"`
	Providers   []ProviderData `json:"providers,omitempty" bson:"providers,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (u *User) IsGranted(roles ...Role) bool {
	if len(roles) == 0 {
		panic(errors.New("roles is empty"))
	}
	return lo.Some(u.Roles, roles)
}

func (u *User) HasPassword() bool {
	return u.Password != ""
}

func (u *User) Secret() string {
	return b32NoPadding.EncodeToString(u.OTPSecret)
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

func (u *User) CheckTOTP(password string) error {
	if totp.Validate(password, u.Secret()) {
		return nil
	}
	return errors.New("otp is not valid")
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
