package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/util"
	"time"
)

const PluginName = "redis"

type Plugin struct {
	cfg *Config
}

func (p *Plugin) Init(cfg config.Configurer) error {
	const op = errors.Op("redis_plugin_init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	if err := cfg.UnmarshalKey(PluginName, &p.cfg); err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*common.RedisMaker)(nil), p.ServiceRedisMaker),
	}
}

func (p *Plugin) ServiceRedisMaker() common.RedisMaker {
	return p
}

func (p *Plugin) Make() (client redis.UniversalClient, err error) {
	client = universalClient(p.cfg)

	if p.cfg.Ping {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		op := errors.Op("redis_maker")

		if res, err1 := client.Ping(ctx).Result(); err1 != nil {
			err = errors.E(op, errors.Errorf("could not check Redis server. %w", err1))
		} else if res != "PONG" {
			err = errors.E(op, "could not check Redis server")
		}
	}
	return
}

func (p *Plugin) MakeRedisClient() any {
	return util.Must(p.Make())
}

func (p *Plugin) Name() string {
	return PluginName
}

func universalClient(cfg *Config) redis.UniversalClient {
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
