package events

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iagapie/rumors/internal/consts"
	"github.com/iagapie/rumors/internal/views"
	"github.com/iagapie/rumors/pkg/emitter"
	"github.com/rs/zerolog"
	"strings"
	"sync"
)

type Listener struct {
	BotAPI  *tgbotapi.BotAPI
	Emitter emitter.Emitter
	Log     zerolog.Logger
	Owner   int64
	done    chan struct{}
}

func (l *Listener) Start() {
	l.done = make(chan struct{}, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go l.run(&wg)
	wg.Wait()
}

func (l *Listener) Shutdown() {
	close(l.done)
}

func (l *Listener) view(chatId int64, template string, data any) {
	text, err := views.View(views.TelegramNS, template, data)
	if err != nil {
		l.Log.Error().Err(err).Str("template", template).Msg("error due build view")
		return
	}
	l.send(chatId, text)
}

func (l *Listener) send(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, "")
	msg.DisableWebPagePreview = true
	msg.ParseMode = "html"

	max := 4096

	var data []string
	for {
		if len(text) > max {
			i := strings.LastIndex(text[:max], "\n")
			if i <= 0 {
				i = max
			}
			data = append(data, text[:i])
			text = text[i+1:]
		} else if len(text) > 0 {
			data = append(data, text)
			break
		} else {
			break
		}
	}

	for _, str := range data {
		msg.Text = str

		if _, err := l.BotAPI.Request(msg); err != nil {
			if e, ok := err.(*tgbotapi.Error); ok {
				l.Log.Error().Interface("error", e).Int64("chat_id", chatId).Msg(text)
			} else {
				l.Log.Error().Err(err).Int64("chat_id", chatId).Msg(text)
			}
			break
		}
	}
}

func (l *Listener) run(wg *sync.WaitGroup) {
	onAppStart := l.Emitter.On(consts.EventAppStart)
	onAppStop := l.Emitter.On(consts.EventAppStop)

	onErrorForbidden := l.Emitter.On(consts.EventErrorForbidden)
	onErrorNotFound := l.Emitter.On(consts.EventErrorNotFound)
	onErrorViewList := l.Emitter.On(consts.EventErrorViewList)
	onErrorArgs := l.Emitter.On(consts.EventErrorArgs)

	onRoomSaveError := l.Emitter.On(consts.EventRoomSaveError)
	onRoomSaveAfter := l.Emitter.On(consts.EventRoomSaveAfter)
	onRoomViewOne := l.Emitter.On(consts.EventRoomViewOne)
	onRoomViewList := l.Emitter.On(consts.EventRoomViewList)

	onFeedSaveError := l.Emitter.On(consts.EventFeedSaveError)
	onFeedSaveAfter := l.Emitter.On(consts.EventFeedSaveAfter)
	onFeedViewOne := l.Emitter.On(consts.EventFeedViewOne)
	onFeedViewList := l.Emitter.On(consts.EventFeedViewList)

	onFeedItemViewList := l.Emitter.On(consts.EventFeedItemViewList)

	defer func(e emitter.Emitter) {
		e.Off(consts.EventAppStart)
		e.Off(consts.EventAppStop)

		e.Off(consts.EventErrorForbidden, onErrorForbidden)
		e.Off(consts.EventErrorNotFound, onErrorNotFound)
		e.Off(consts.EventErrorViewList, onErrorViewList)
		e.Off(consts.EventErrorArgs, onErrorArgs)

		e.Off(consts.EventRoomSaveError, onRoomSaveError)
		e.Off(consts.EventRoomSaveAfter, onRoomSaveAfter)
		e.Off(consts.EventRoomViewOne, onRoomViewOne)
		e.Off(consts.EventRoomViewList, onRoomViewList)

		e.Off(consts.EventFeedViewList, onFeedViewList)
		e.Off(consts.EventFeedViewOne, onFeedViewOne)
		e.Off(consts.EventFeedSaveAfter, onFeedSaveAfter)
		e.Off(consts.EventFeedSaveError, onFeedSaveError)

		e.Off(consts.EventFeedItemViewList, onFeedItemViewList)
	}(l.Emitter)

	wg.Done()

	for {
		select {
		case <-l.done:
			return

		case event := <-onAppStart:
			l.onAppStart(event)
		case event := <-onAppStop:
			l.onAppStop(event)

		case event := <-onErrorForbidden:
			l.onErrorForbidden(event)
		case event := <-onErrorNotFound:
			l.onErrorNotFound(event)
		case event := <-onErrorViewList:
			l.onErrorViewList(event)
		case event := <-onErrorArgs:
			l.onErrorArgs(event)

		case event := <-onRoomSaveError:
			l.onRoomSaveError(event)
		case event := <-onRoomSaveAfter:
			l.onRoomSaveAfter(event)
		case event := <-onRoomViewOne:
			l.onRoomViewOne(event)
		case event := <-onRoomViewList:
			l.onRoomViewList(event)

		case event := <-onFeedSaveError:
			l.onFeedSaveError(event)
		case event := <-onFeedSaveAfter:
			l.onFeedSaveAfter(event)
		case event := <-onFeedViewOne:
			l.onFeedViewOne(event)
		case event := <-onFeedViewList:
			l.onFeedViewList(event)

		case event := <-onFeedItemViewList:
			l.onFeedItemViewList(event)
		}
	}
}
