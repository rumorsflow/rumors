package events

import "github.com/olebedev/emitter"

func (l *Listener) onErrorForbidden(event emitter.Event) {
	l.view(event.Args[0].(int64), "forbidden.html", nil)
}

func (l *Listener) onErrorNotFound(event emitter.Event) {
	l.view(event.Args[0].(int64), "notfound.html", event.Args[1].(error).Error())
}

func (l *Listener) onErrorViewList(event emitter.Event) {
	l.view(event.Args[0].(int64), "error.html", event.Args[1].(error).Error())
}

func (l *Listener) onErrorArgs(event emitter.Event) {
	l.view(event.Args[0].(int64), "error.html", event.Args[1].(string))
}
