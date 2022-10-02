package events

import (
	"github.com/iagapie/rumors/internal/models"
	"github.com/olebedev/emitter"
)

func (l *Listener) onFeedSaveError(event emitter.Event) {
	l.view(l.Owner, "error.html", event.Args[1].(error).Error())
	l.onFeedSaveAfter(event)
}

func (l *Listener) onFeedSaveAfter(event emitter.Event) {
	l.view(l.Owner, "feed.html", event.Args[0].(models.Feed))
}

func (l *Listener) onFeedViewOne(event emitter.Event) {
	l.view(event.Args[0].(int64), "feed.html", event.Args[1].(models.Feed))
}

func (l *Listener) onFeedViewList(event emitter.Event) {
	l.view(event.Args[0].(int64), "feeds.html", event.Args[1].([]models.Feed))
}
