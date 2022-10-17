package tgbotsender

import (
	"bytes"
	"embed"
	"github.com/go-fc/slice"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/pkg/str"
	"go.uber.org/zap"
	"html/template"
	"io"
	"strings"
)

type View string

const (
	ViewRoom      = "room.html"
	ViewFeedItems = "feeditems.html"
	ViewAppStart  = "appstart.html"
	ViewAppStop   = "appstop.html"
)

type TelegramSender interface {
	Owner() int64
	SendView(chatId int64, view View, data any)
	SendText(chatId int64, text string)
}

var (
	//go:embed views/*.html
	viewsFS embed.FS

	replacer = strings.NewReplacer(".", "", "-", "", " ", "")
	funcMap  = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"join": func(data any, sep string) string {
			if data == nil {
				return ""
			}

			switch tmp := data.(type) {
			case *[]string:
				return strings.Join(*tmp, sep)
			case []string:
				return strings.Join(tmp, sep)
			case *[]models.RoomPermission:
				return strings.Join(slice.Map(*tmp, func(p models.RoomPermission) string {
					return string(p)
				}), sep)
			}
			return ""
		},
		"hashtag": func(tags []string) string {
			var tmp []string
			for _, tag := range tags {
				if tag = replacer.Replace(tag); tag != "" {
					tmp = append(tmp, "#"+tag)
				}

			}
			return strings.Join(tmp, " ")
		},
	}
)

func (p *Plugin) Owner() int64 {
	return p.cfg.Owner
}

func (p *Plugin) SendView(chatId int64, view View, data any) {
	p.log.Debug("bot send view", zap.Int64("chatId", chatId), zap.Any("view", view), zap.Any("data", data))

	text, err := p.view(view, data)
	if err != nil {
		p.log.Error("error due build view", zap.Any("view", view), zap.Error(err))
		return
	}
	p.SendText(chatId, text)
}

func (p *Plugin) SendText(chatId int64, text string) {
	if chatId == 0 {
		chatId = p.cfg.Owner
	}

	if chatId == 0 {
		p.log.Warn("error due the chat id is zero")
		return
	}

	msg := tgbotapi.NewMessage(chatId, text)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "html"

	p.log.Debug("bot send message", zap.Any("message", msg))

	for _, chunk := range str.SplitMax(text, "\n", 4096) {
		msg.Text = chunk

		if _, err := p.botApi.Request(msg); err != nil {
			if e, ok := err.(*tgbotapi.Error); ok {
				p.log.Error(text, zap.Int64("chatId", chatId), zap.Error(e), zap.Int("errorCode", e.Code))
			} else {
				p.log.Error(text, zap.Int64("chatId", chatId), zap.Error(e))
			}
			break
		}
	}
}

func (p *Plugin) initTemplates() (err error) {
	p.templates, err = template.
		New("telegram").
		Funcs(funcMap).
		ParseFS(viewsFS, "views/*")
	return
}

func (p *Plugin) execute(w io.Writer, view View, data any) error {
	p.Lock()
	defer p.Unlock()

	return p.templates.ExecuteTemplate(w, string(view), data)
}

func (p *Plugin) view(view View, data any) (string, error) {
	var out bytes.Buffer
	if err := p.execute(&out, view, data); err != nil {
		return "", err
	}
	return strings.Trim(out.String(), "\n"), nil
}
