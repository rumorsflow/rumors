package entity

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

type SysUser struct {
	ID        uuid.UUID `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string    `json:"username,omitempty" bson:"username,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Password  string    `json:"-" bson:"password,omitempty"`
	OTPSecret []byte    `json:"-" bson:"otp_secret,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (u *SysUser) EntityID() uuid.UUID {
	return u.ID
}

func (u *SysUser) Secret() string {
	return b32NoPadding.EncodeToString(u.OTPSecret)
}

func (u *SysUser) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}

func (u *SysUser) CheckTOTP(password string) error {
	if totp.Validate(password, u.Secret()) {
		return nil
	}
	return errors.New("otp is not valid")
}

func (u *SysUser) GenerateOTPSecret(size uint) (err error) {
	u.OTPSecret, err = generateSecret(size)
	return err
}

func (u *SysUser) GeneratePasswordHash() error {
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
	if size == 0 {
		size = 20
	}

	s := make([]byte, size)
	_, err := rand.Reader.Read(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
