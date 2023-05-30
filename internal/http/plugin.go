package http

import (
	"context"
	"github.com/gowool/middleware/cors"
	"github.com/gowool/middleware/gzip"
	"github.com/gowool/middleware/prometheus"
	"github.com/gowool/middleware/proxy"
	"github.com/gowool/middleware/sse"
	"github.com/gowool/swagger"
	"github.com/gowool/wool"
	"github.com/redis/go-redis/v9"
	"github.com/roadrunner-server/errors"
	_ "github.com/rumorsflow/rumors/v2/docs"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/front"
	"github.com/rumorsflow/rumors/v2/internal/http/sys"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/jwt"
	"github.com/rumorsflow/rumors/v2/pkg/logger"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/pprof"
)

const PluginName = "http"

type Plugin struct {
	client       redis.UniversalClient
	queueActions *sys.QueueActions
	srv          *wool.Server
	w            *wool.Wool
	front        *front.Front
	sys          *sys.Sys
	done         chan struct{}
}

func (p *Plugin) Init(cfg config.Configurer, rdbMaker common.RedisMaker, sub common.Sub, uow common.UnitOfWork, log logger.Logger) error {
	const op = errors.Op("http_plugin_init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	var httpCfg Config
	if err := cfg.UnmarshalKey(PluginName, &httpCfg); err != nil {
		return errors.E(op, err)
	}

	var srvCfg wool.ServerConfig
	if err := cfg.UnmarshalKey(PluginName, &srvCfg); err != nil {
		return errors.E(op, err)
	}

	httpCfg.init(cfg.Version())
	srvCfg.Init()

	siteAny, err := uow.Repository((*entity.Site)(nil))
	if err != nil {
		return errors.E(op, err)
	}
	siteRepo := siteAny.(repository.ReadWriteRepository[*entity.Site])

	articleAny, err := uow.Repository((*entity.Article)(nil))
	if err != nil {
		return errors.E(op, err)
	}
	articleRepo := articleAny.(repository.ReadWriteRepository[*entity.Article])

	chatAny, err := uow.Repository((*entity.Chat)(nil))
	if err != nil {
		return errors.E(op, err)
	}
	chatRepo := chatAny.(repository.ReadWriteRepository[*entity.Chat])

	jobAny, err := uow.Repository((*entity.Job)(nil))
	if err != nil {
		return errors.E(op, err)
	}
	jobRepo := jobAny.(repository.ReadWriteRepository[*entity.Job])

	sysUserAny, err := uow.Repository((*entity.SysUser)(nil))
	if err != nil {
		return errors.E(op, err)
	}
	sysUserRepo := sysUserAny.(repository.ReadWriteRepository[*entity.SysUser])

	client, err := rdbMaker.Make()
	if err != nil {
		return errors.E(op, err)
	}

	signer := jwt.NewSigner(httpCfg.JWT.GetPrivateKey())
	authService := sys.NewAuthService(sysUserRepo, client, signer, httpCfg.JWT)
	p.queueActions = sys.NewQueueActions(rdbMaker)
	p.client = client

	l := log.NamedLogger(PluginName)
	frontLog := l.WithGroup("front")
	sysLog := l.WithGroup("sys")

	p.srv = wool.NewServer(&srvCfg, l.WithGroup("server"))
	p.w = wool.New(
		l,
		wool.WithErrorTransform(ErrorTransform),
		wool.WithAfterServe(AfterServe(httpCfg.LogReq, l)),
	)

	if httpCfg.Middleware.Metrics != nil {
		p.w.Use(prometheus.Middleware(httpCfg.Middleware.Metrics))
	}

	p.w.Use(
		proxy.Middleware(),
		gzip.Middleware(httpCfg.Middleware.Compress),
		cors.Middleware(httpCfg.Middleware.CORS),
	)

	p.w.MountHealth()

	prometheus.Mount(p.w)

	if l.Enabled(context.Background(), slog.LevelDebug) {
		p.w.Group("/debug/pprof", func(pp *wool.Wool) {
			pp.Add("/cmdline", wool.ToHandler(http.HandlerFunc(pprof.Cmdline)))
			pp.Add("/profile", wool.ToHandler(http.HandlerFunc(pprof.Profile)))
			pp.Add("/symbol", wool.ToHandler(http.HandlerFunc(pprof.Symbol)))
			pp.Add("/trace", wool.ToHandler(http.HandlerFunc(pprof.Trace)))
			pp.Add("/...", wool.ToHandler(http.HandlerFunc(pprof.Index)))
		})
	}

	if httpCfg.Swagger.Enabled {
		p.w.Group("/swagger", func(sw *wool.Wool) {
			sw.GET("/sys/...", swagger.New(httpCfg.Swagger.Sys).Handler)
			sw.GET("/front/...", swagger.New(httpCfg.Swagger.Front).Handler)
		})
	}

	p.sys = &sys.Sys{
		Logger:         sysLog,
		CfgJWT:         httpCfg.JWT,
		DirUI:          httpCfg.UI.SysPath,
		QueueActions:   p.queueActions,
		SSE:            sys.NewSSE(rdbMaker, sysLog.WithGroup("sse")),
		AuthActions:    sys.NewAuthActions(authService, sysLog.WithGroup("auth")),
		ArticleActions: sys.NewArticleActions(articleRepo, articleRepo),
		SiteCRUD:       sys.NewSiteCRUD(siteRepo, siteRepo),
		ChatCRUD:       sys.NewChatCRUD(chatRepo, chatRepo),
		JobCRUD:        sys.NewJobCRUD(jobRepo, jobRepo),
	}

	p.front = &front.Front{
		Logger:         frontLog,
		Sub:            sub,
		DirUI:          httpCfg.UI.FrontPath,
		SSE:            sse.New(httpCfg.SSE, frontLog.WithGroup("sse")),
		SiteActions:    &front.SiteActions{SiteRepo: siteRepo},
		ArticleActions: &front.ArticleActions{ArticleRepo: articleRepo, SiteRepo: siteRepo},
	}

	p.w.Group("", func(sw *wool.Wool) {
		p.sys.Register(sw)
	})

	p.w.Group("", func(fw *wool.Wool) {
		p.front.Register(fw)
	})

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	go func(w *wool.Wool, srv *wool.Server, errCh chan<- error) {
		if err := srv.Start(w); err != nil {
			errCh <- err
		}
	}(p.w, p.srv, errCh)

	p.done = make(chan struct{})
	p.front.Listen(p.done)
	p.sys.Listen(p.done)

	return errCh
}

func (p *Plugin) Stop(ctx context.Context) error {
	close(p.done)

	err := p.srv.Shutdown(ctx)
	err = errs.Append(err, p.queueActions.Close())
	err = errs.Append(err, p.client.Close())
	return err
}

func (p *Plugin) Name() string {
	return PluginName
}
