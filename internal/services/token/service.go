package token

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/pkg/jwt"
)

type Session struct {
	ExpiresIn    uint   `json:"expires_in"` // The lifetime in seconds of refresh token
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Service interface {
	ValidateJWT(ctx context.Context, token string) (jwt.UserClaims, models.User, error)
	SessByRefreshToken(ctx context.Context, refreshToken string) (Session, error)
	SessByUser(ctx context.Context, user models.User, metadata any) (Session, error)
}

func (p *Plugin) ValidateJWT(ctx context.Context, token string) (jwt.UserClaims, models.User, error) {
	_, claims, err := p.verifier.Parse(token)
	if err != nil {
		return claims, models.User{}, fmt.Errorf("invalid access token: %w", err)
	}

	user, err := p.userStorage.FindById(ctx, claims.ID)
	if err != nil {
		err = fmt.Errorf("invalid access token: %w", err)
	}
	return claims, user, err
}

func (p *Plugin) SessByRefreshToken(ctx context.Context, refreshToken string) (Session, error) {
	data, err := p.storage.Get(ctx, refreshToken)
	if err != nil {
		return Session{}, fmt.Errorf("invalid refresh token: %w", err)
	}
	defer p.storage.Del(ctx, refreshToken)

	user, err := p.userStorage.FindById(ctx, data.Id)
	if err != nil {
		return Session{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	return p.SessByUser(ctx, user, data.Metadata)
}

func (p *Plugin) SessByUser(ctx context.Context, user models.User, metadata any) (Session, error) {
	refreshToken := uuid.New().String()
	if err := p.storage.Set(ctx, refreshToken, StorageData{Id: user.Id, Metadata: metadata}, p.cfg.TTL.Refresh); err != nil {
		return Session{}, fmt.Errorf("save refresh token error: %w", err)
	}
	claims := jwt.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims(user.Id, p.cfg.TTL.JWT),
		Username:         user.Username,
		Email:            user.Email,
		Roles:            user.Roles,
		Providers:        user.Providers,
		Metadata:         metadata,
	}
	accessToken, err := p.signer.Sign(claims)
	if err != nil {
		return Session{}, fmt.Errorf("generate access token error: %w", err)
	}
	return Session{
		ExpiresIn:    uint(p.cfg.TTL.Refresh.Seconds()),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
