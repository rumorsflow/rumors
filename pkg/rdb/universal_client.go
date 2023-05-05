package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
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
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:                 cfg.Addrs,
		ClientName:            cfg.ClientName,
		DB:                    cfg.DB,
		Username:              cfg.Username,
		Password:              cfg.Password,
		SentinelUsername:      cfg.SentinelUsername,
		SentinelPassword:      cfg.SentinelPassword,
		MaxRetries:            cfg.MaxRetries,
		MinRetryBackoff:       cfg.MaxRetryBackoff,
		MaxRetryBackoff:       cfg.MaxRetryBackoff,
		DialTimeout:           cfg.DialTimeout,
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		ContextTimeoutEnabled: cfg.ContextTimeoutEnabled,
		PoolFIFO:              cfg.PoolFIFO,
		PoolSize:              cfg.PoolSize,
		PoolTimeout:           cfg.PoolTimeout,
		MinIdleConns:          cfg.MinIdleConns,
		MaxIdleConns:          cfg.MaxIdleConns,
		ConnMaxIdleTime:       cfg.ConnMaxIdleTime,
		ConnMaxLifetime:       cfg.ConnMaxLifetime,
		MaxRedirects:          cfg.MaxRedirects,
		ReadOnly:              cfg.ReadOnly,
		RouteByLatency:        cfg.RouteByLatency,
		RouteRandomly:         cfg.RouteRandomly,
		MasterName:            cfg.MasterName,
	})
}
