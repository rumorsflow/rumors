package front

import (
	"context"
	"github.com/gowool/middleware/sse"
	"github.com/gowool/wool"
	"github.com/gowool/wool/render"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"golang.org/x/exp/slog"
)

var uiBuiltIn = true

type Front struct {
	Logger         *slog.Logger
	Sub            *pubsub.Subscriber
	SSE            *sse.Event
	FeedActions    *FeedActions
	ArticleActions *ArticleActions
}

func (front *Front) Register(mux *wool.Wool) {
	mux.Group("/api/v1", func(w *wool.Wool) {
		w.GET("/feeds", front.FeedActions.List)
		w.GET("/articles", front.ArticleActions.List)

		w.Group("", func(sw *wool.Wool) {
			sw.Use(front.SSE.Middleware)
			sw.GET("/realtime", front.SSE.Handler)
		})
	})

	front.Logger.WithGroup("api").WithGroup("v1").Info("frontend V1 APIs registered")

	if uiBuiltIn {
		mux.UI("", assetFS())

		front.Logger.WithGroup("ui").Info("frontend UI registered")
	}
}

func (front *Front) Listen(ctx context.Context) error {
	articlesCh := front.Sub.Articles(ctx).Channel()

	for {
		select {
		case <-ctx.Done():
			return front.SSE.Close()
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
