package rdb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"time"
)

type UniversalClientMaker struct {
	c di.Container
}

func (m *UniversalClientMaker) Make() redis.UniversalClient {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	return di.Must(NewUniversalClient(ctx, m.c))
}

func (m *UniversalClientMaker) MakeRedisClient() any {
	return m.Make()
}

func New(cfg *Config) redis.UniversalClient {
	cfg.Init()

	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              cfg.Addrs,
		DB:                 cfg.DB,
		Username:           cfg.Username,
		Password:           cfg.Password,
		SentinelPassword:   cfg.SentinelPassword,
		MaxRetries:         cfg.MaxRetries,
		MinRetryBackoff:    cfg.MaxRetryBackoff,
		MaxRetryBackoff:    cfg.MaxRetryBackoff,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		PoolSize:           cfg.PoolSize,
		MinIdleConns:       cfg.MinIdleConns,
		MaxConnAge:         cfg.MaxConnAge,
		PoolTimeout:        cfg.PoolTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFreq,
		ReadOnly:           cfg.ReadOnly,
		RouteByLatency:     cfg.RouteByLatency,
		RouteRandomly:      cfg.RouteRandomly,
		MasterName:         cfg.MasterName,
	})
}
