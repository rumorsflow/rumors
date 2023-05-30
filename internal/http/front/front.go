package front

import (
	"context"
	"github.com/gowool/middleware/cfipcountry"
	"github.com/gowool/middleware/sse"
	"github.com/gowool/wool"
	"github.com/gowool/wool/render"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"golang.org/x/exp/slog"
	"net/http"
)

var uiBuiltIn = true

type Front struct {
	Logger         *slog.Logger
	Sub            common.Sub
	SSE            *sse.Event
	SiteActions    *SiteActions
	ArticleActions *ArticleActions
	CfConfig       *cfipcountry.Config
	DirUI          string
}

func (front *Front) Register(mux *wool.Wool) {
	mux.Group("/api/v1", func(w *wool.Wool) {
		w.GET("/sites", front.SiteActions.List)
		w.GET("/articles", front.ArticleActions.List)

		w.Group("", func(sw *wool.Wool) {
			sw.Use(front.SSE.Middleware)
			sw.GET("/realtime", front.SSE.Handler)
		})
	})

	front.Logger.WithGroup("api").WithGroup("v1").Info("frontend V1 APIs registered")

	if front.CfConfig != nil {
		mux.Use(cfipcountry.Middleware(front.CfConfig))
	}

	if front.DirUI != "" {
		mux.UI("", http.Dir(front.DirUI))
		front.Logger.WithGroup("ui").Info("frontend UI registered")
	} else if uiBuiltIn {
		mux.UI("", assetFS())
		front.Logger.WithGroup("ui").Info("frontend UI registered")
	}
}

func (front *Front) Listen(done <-chan struct{}) {
	go front.listen(done)
}

func (front *Front) listen(done <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
		_ = front.SSE.Close()
	}()

	articlesCh := front.Sub.Articles(ctx).Channel()

	for {
		select {
		case <-done:
			return
		case data := <-articlesCh:
			if data == nil {
				continue
			}
			front.SSE.Broadcast(render.SSEvent{
				Event: "articles",
				Data:  data.Payload,
			})
		}
	}
}
