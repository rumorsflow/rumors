package sys

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/conv"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
)

var _ jwtv4.Claims = (*UserClaims)(nil)

const (
	issuer        = "rumors-sys-api"
	audience      = "rumors-sys"
	frontAudience = "rumors-front"

	filter = "index=0&size=1&field.0.0=username&value.0.0=%[1]s&field.0.1=email&value.0.1=%[1]s"
)

type UserClaims struct {
	jwt.UserClaims
}

func (c UserClaims) Valid() error {
	vErr := &jwtv4.ValidationError{}
	if err := c.UserClaims.Valid(); err != nil {
		vErr = err.(*jwtv4.ValidationError)
	}

	if !c.VerifyIssuer(issuer, true) {
		vErr.Inner = errs.Append(vErr.Inner, jwtv4.ErrTokenInvalidIssuer)
		vErr.Errors |= jwtv4.ValidationErrorIssuer
	}

	if !c.VerifyAudience(audience, true) {
		vErr.Inner = errs.Append(vErr.Inner, jwtv4.ErrTokenInvalidAudience)
		vErr.Errors |= jwtv4.ValidationErrorAudience
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

type Session struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type redisData struct {
	Username string `json:"username"`
	OTP      bool   `json:"otp,omitempty"`
}

type AuthService interface {
	SignInByCredentials(ctx context.Context, username, password string) (Session, error)
	SignInByOTP(ctx context.Context, username, password string) (Session, error)
	SignInByRefreshToken(ctx context.Context, refreshToken string) (Session, error)
}

type authService struct {
	userRepo repository.ReadRepository[*entity.SysUser]
	client   redis.UniversalClient
	signer   jwt.Signer
	cfgJWT   *jwt.Config
}

func NewAuthService(
	userRepo repository.ReadRepository[*entity.SysUser],
	client redis.UniversalClient,
	signer jwt.Signer,
	cfgJWT *jwt.Config,
) AuthService {
	cfgJWT.Init()

	return &authService{
		userRepo: userRepo,
		client:   client,
		signer:   signer,
		cfgJWT:   cfgJWT,
	}
}

func (s *authService) SignInByCredentials(ctx context.Context, username, password string) (Session, error) {
	user, err := s.findUserBy(ctx, username)
	if err != nil {
		return Session{}, err
	}
	if err = user.CheckPassword(password); err != nil {
		return Session{}, err
	}

	return s.sessionByUser(ctx, user, false)
}

func (s *authService) SignInByOTP(ctx context.Context, username, password string) (Session, error) {
	user, err := s.findUserBy(ctx, username)
	if err != nil {
		return Session{}, err
	}
	if err = user.CheckTOTP(password); err != nil {
		return Session{}, err
	}

	return s.sessionByUser(ctx, user, true)
}

func (s *authService) SignInByRefreshToken(ctx context.Context, refreshToken string) (Session, error) {
	raw, err := s.client.Get(ctx, refreshToken).Result()
	if err != nil {
		return Session{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	defer s.client.Del(ctx, refreshToken)

	var data redisData
	if err = json.Unmarshal(conv.StringToBytes(raw), &data); err != nil {
		return Session{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.findUserBy(ctx, data.Username)
	if err != nil {
		return Session{}, err
	}

	return s.sessionByUser(ctx, user, data.OTP)
}

func (s *authService) sessionByUser(ctx context.Context, user *entity.SysUser, otp bool) (Session, error) {
	claims := UserClaims{
		UserClaims: jwt.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims(
				issuer, user.ID,
				jwtv4.ClaimStrings{audience, frontAudience},
				s.cfgJWT.AccessTokenTTL,
			),
			Username: user.Username,
			Email:    user.Email,
			OTP:      otp,
		},
	}

	accessToken, err := s.signer.Sign(claims)
	if err != nil {
		return Session{}, err
	}

	data, err := json.Marshal(redisData{Username: user.Username, OTP: otp})
	if err != nil {
		return Session{}, err
	}

	refreshToken := "dummy"

	if otp {
		refreshToken = uuid.NewString()
		if err = s.client.Set(ctx, refreshToken, conv.BytesToString(data), s.cfgJWT.RefreshTokenTTL).Err(); err != nil {
			return Session{}, fmt.Errorf("save refresh token error: %w", err)
		}
	}

	return Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) findUserBy(ctx context.Context, username string) (*entity.SysUser, error) {
	criteria := db.BuildCriteria(fmt.Sprintf(filter, username))
	users, err := s.userRepo.Find(ctx, criteria)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, repository.ErrEntityNotFound
	}

	return users[0], nil
}
