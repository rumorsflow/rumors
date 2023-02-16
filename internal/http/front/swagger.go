//go:build nop

package front

//	@title		Rumors Frontend API
//	@version	1.0
//	@BasePath	/api/v1

//	@Summary		List feeds
//	@Description	get feeds
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int			false	"Page Index"	default(0)	minimum(0)
//	@Param			size	query		int			false	"Page Size"		default(20)	minimum(1)	maximum(100)
//	@Success		200		{array}		entity.Feed	"OK"
//	@Failure		400		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/feeds [get]
func nopFeedList() {}

//	@Summary		List articles
//	@Description	get articles
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int				false	"Page Index"	default(0)	minimum(0)
//	@Param			size	query		int				false	"Page Size"		default(20)	minimum(1)	maximum(100)
//	@Param			h		query		string			false	"Source Hosts"
//	@Param			dt		query		string			false	"From DateTime"	Format(date-time)
//	@Param			l		query		string			false	"Languages"
//	@Success		200		{array}		pubsub.Article	"OK"
//	@Failure		400		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/articles [get]
func nopArticleList() {}

//	@Summary		Realtime
//	@Description	sse stream
//	@Tags			sse
//	@Header			default	{string}	Content-Type	text/event-stream
//	@response		default
//	@Failure		400	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/realtime [get]
func nopRealtime() {}
