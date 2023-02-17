package http

import (
	"context"
	"crypto/tls"
	"github.com/gowool/middleware/cors"
	"github.com/gowool/middleware/gzip"
	"github.com/gowool/middleware/prometheus"
	"github.com/gowool/middleware/proxy"
	"github.com/gowool/middleware/sse"
	"github.com/gowool/wool"
	_ "github.com/rumorsflow/rumors/v2/docs"
	"github.com/rumorsflow/rumors/v2/internal/http/front"
	"github.com/rumorsflow/rumors/v2/internal/http/sys"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/rdb"
	"io/fs"
	"net/http"
	"net/http/pprof"
	"time"
)

const (
	ConfigServerKey   = "http"
	ConfigCORSKey     = "http.middleware.cors"
	ConfigMetricsKey  = "http.middleware.metrics"
	ConfigCompressKey = "http.middleware.compress"
)

type (
	WoolKey   struct{}
	ServerKey struct{}
	FrontKey  struct{}
	SysKey    struct{}
)

func GetWool(ctx context.Context, c ...di.Container) (*wool.Wool, error) {
	return di.Get[*wool.Wool](ctx, WoolKey{}, c...)
}

func GetServer(ctx context.Context, c ...di.Container) (*wool.Server, error) {
	return di.Get[*wool.Server](ctx, ServerKey{}, c...)
}

func FrontActivator(version string) *di.Activator {
	return &di.Activator{
		Key: FrontKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			siteRepo, err := db.GetSiteRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			articleRepo, err := db.GetArticleRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			sub, err := pubsub.GetSub(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			sseCfg := &sse.Config{Version: version, ClientIdle: time.Hour}

			log := logger.WithGroup("http").WithGroup("front")

			return &front.Front{
				Logger:         log,
				Sub:            sub,
				SSE:            sse.New(sseCfg, log.WithGroup("sse")),
				SiteActions:    &front.SiteActions{SiteRepo: siteRepo},
				ArticleActions: &front.ArticleActions{ArticleRepo: articleRepo, SiteRepo: siteRepo},
			}, nil, nil
		}),
	}
}

func SysActivator() *di.Activator {
	return &di.Activator{
		Key: SysKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			siteRepo, err := db.GetSiteRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			feedRepo, err := db.GetFeedRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			chatRepo, err := db.GetChatRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			jobRepo, err := db.GetJobRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			articleRepo, err := db.GetArticleRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			userRepo, err := db.GetSysUserRepository(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			jwtCfg, err := jwt.GetConfig(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			signer, err := jwt.GetSigner(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			client, err := rdb.NewUniversalClient(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			cl := di.CloserFunc(func(context.Context) error {
				return client.Close()
			})

			authService := sys.NewAuthService(userRepo, client, signer, jwtCfg)

			return &sys.Sys{
				Logger:         logger.WithGroup("http").WithGroup("sys"),
				CfgJWT:         jwtCfg,
				AuthActions:    sys.NewAuthActions(authService),
				ArticleActions: sys.NewArticleActions(articleRepo, articleRepo),
				SiteCRUD:       sys.NewSiteCRUD(siteRepo, siteRepo),
				FeedCRUD:       sys.NewFeedCRUD(feedRepo, feedRepo),
				ChatCRUD:       sys.NewChatCRUD(chatRepo, chatRepo),
				JobCRUD:        sys.NewJobCRUD(jobRepo, jobRepo),
			}, cl, nil
		}),
	}
}

func WoolActivator(version string) *di.Activator {
	return &di.Activator{
		Key: WoolKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			compressConfig, err := config.UnmarshalKey[*gzip.Config](c.Configurer(), ConfigCompressKey)
			if err != nil {
				return nil, nil, errs.E(di.OpFactory, err)
			}

			corsConfig, err := config.UnmarshalKey[*cors.Config](c.Configurer(), ConfigCORSKey)
			if err != nil {
				return nil, nil, errs.E(di.OpFactory, err)
			}

			wool.SetLogger(logger.WithGroup("http"))

			w := wool.New(
				wool.WithErrorTransform(ErrorTransform),
				wool.WithAfterServe(AfterServe(nil)),
			)

			if metricsConfig, err := config.UnmarshalKeyE[*prometheus.Config](c.Configurer(), ConfigMetricsKey); err == nil {
				metricsConfig.Version = version
				w.Use(prometheus.Middleware(metricsConfig))
			}

			w.Use(
				proxy.Middleware(),
				gzip.Middleware(compressConfig),
				cors.Middleware(corsConfig),
			)

			if logger.IsDebug() {
				w.Group("/debug/pprof", func(pp *wool.Wool) {
					pp.Add("/cmdline", wool.ToHandler(http.HandlerFunc(pprof.Cmdline)))
					pp.Add("/profile", wool.ToHandler(http.HandlerFunc(pprof.Profile)))
					pp.Add("/symbol", wool.ToHandler(http.HandlerFunc(pprof.Symbol)))
					pp.Add("/trace", wool.ToHandler(http.HandlerFunc(pprof.Trace)))
					pp.Add("/...", wool.ToHandler(http.HandlerFunc(pprof.Index)))
				})
			}

			return w, nil, nil
		}),
	}
}

func ServerActivator(certFS fs.FS, tls func(*tls.Config)) *di.Activator {
	return &di.Activator{
		Key: ServerKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			cfg, err := config.UnmarshalKey[*wool.ServerConfig](c.Configurer(), ConfigServerKey)
			if err != nil {
				return nil, nil, errs.E(di.OpFactory, err)
			}

			s := wool.NewServer(cfg)
			s.CertFilesystem = certFS
			s.TLSConfig = tls

			return s, nil, nil
		}),
	}
}
