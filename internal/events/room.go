package events

import (
	"github.com/iagapie/rumors/internal/models"
	"github.com/olebedev/emitter"
)

func (l *Listener) onRoomSaveError(event emitter.Event) {
	l.view(l.Owner, "error.html", event.Args[1].(error).Error())
	l.onRoomSaveAfter(event)
}

func (l *Listener) onRoomSaveAfter(event emitter.Event) {
	room := event.Args[0].(models.Room)
	if room.ChatId != l.Owner {
		l.view(l.Owner, "room.html", event.Args[0].(models.Room))
	}
}

func (l *Listener) onRoomViewOne(event emitter.Event) {
	l.view(event.Args[0].(int64), "room.html", event.Args[1].(models.Room))
}

func (l *Listener) onRoomViewList(event emitter.Event) {
	l.view(event.Args[0].(int64), "rooms.html", event.Args[1].([]models.Room))
}
