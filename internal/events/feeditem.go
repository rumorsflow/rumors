package events

import (
	"github.com/iagapie/rumors/internal/models"
	"github.com/olebedev/emitter"
)

func (l *Listener) onFeedItemViewList(event emitter.Event) {
	l.view(event.Args[0].(int64), "feeditems.html", event.Args[1].(map[string][]models.FeedItem))
}
