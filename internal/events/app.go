package events

import "github.com/olebedev/emitter"

func (l *Listener) onAppStart(_ emitter.Event) {
	l.view(l.Owner, "appstart.html", nil)
}

func (l *Listener) onAppStop(_ emitter.Event) {
	l.view(l.Owner, "appstop.html", nil)
}
