package token

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

var _ Storage = (*tokenStorage)(nil)

type StorageData struct {
	Id       string `json:"id,omitempty"`
	Metadata any    `json:"metadata,omitempty"`
}

type Storage interface {
	Set(ctx context.Context, refreshToken string, data StorageData, expiration time.Duration) error
	Get(ctx context.Context, refreshToken string) (StorageData, error)
	Del(ctx context.Context, refreshToken string)
}

type tokenStorage struct {
	client redis.UniversalClient
}

func (s *tokenStorage) Set(ctx context.Context, refreshToken string, data StorageData, expiration time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, refreshToken, string(bytes), expiration).Err()
}

func (s *tokenStorage) Get(ctx context.Context, refreshToken string) (data StorageData, err error) {
	str, err := s.client.Get(ctx, refreshToken).Result()
	if err == nil {
		err = json.Unmarshal([]byte(str), &data)
	}
	return
}

func (s *tokenStorage) Del(ctx context.Context, refreshToken string) {
	s.client.Del(ctx, refreshToken)
}
