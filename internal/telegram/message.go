package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goccy/go-json"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/strutil"
)

type View string

const (
	ViewAppStart View = "appstart.html"
	ViewAppStop  View = "appstop.html"
	ViewArticles View = "articles.html"
	ViewArticle  View = "article.html"
	ViewChat     View = "chat.html"
	ViewSites    View = "sites.html"
	ViewSub      View = "sub.html"
	ViewSuccess  View = "success.html"
	ViewError    View = "error.html"
	ViewNotFound View = "notfound.html"

	OpBotNew  errs.Op = "bot: new"
	OpBotSend errs.Op = "bot: send"
)

type Message struct {
	ChatID   int64  `json:"chat_id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	View     View   `json:"view,omitempty"`
	Data     any    `json:"data,omitempty"`
	Delay    bool   `json:"delay,omitempty"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var msg struct {
		ChatID   int64           `json:"chat_id,omitempty"`
		ImageURL string          `json:"image_url,omitempty"`
		View     View            `json:"view,omitempty"`
		Data     json.RawMessage `json:"data,omitempty"`
		Delay    bool            `json:"delay,omitempty"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	m.ChatID = msg.ChatID
	m.ImageURL = msg.ImageURL
	m.View = msg.View
	m.Delay = msg.Delay

	if len(msg.Data) == 0 {
		return nil
	}

	switch msg.View {
	case ViewArticles:
		var tmp map[string][]pubsub.Article
		if err := json.Unmarshal(msg.Data, &tmp); err != nil {
			return err
		}
		m.Data = tmp
	case ViewChat:
		var tmp entity.Chat
		if err := json.Unmarshal(msg.Data, &tmp); err != nil {
			return err
		}
		m.Data = tmp
	case ViewSites, ViewSub:
		var tmp []string
		if err := json.Unmarshal(msg.Data, &tmp); err != nil {
			return err
		}
		m.Data = tmp
	default:
		var tmp string
		if err := json.Unmarshal(msg.Data, &tmp); err != nil {
			return err
		}
		m.Data = tmp
	}

	return nil
}

func (m *Message) unmarshalData(data json.RawMessage, i any) error {
	if err := json.Unmarshal(data, i); err != nil {
		return err
	}
	m.Data = i
	return nil
}

func (m *Message) chattable(bot *Bot) ([]tgbotapi.Chattable, error) {
	if m.ChatID == 0 {
		m.ChatID = bot.cfg.OwnerID
	}

	if m.View == "" {
		return nil, errs.E(OpBotSend, "message view is required")
	}

	text, err := view(m.View, m.Data)
	if err != nil {
		return nil, errs.E(OpBotSend, "error due to execute template", err)
	}

	chunks := strutil.SplitMax(text, "\n", 4096)
	if len(chunks) == 0 {
		return nil, errs.E(OpBotSend, "error due to split text in chunks")
	}

	messages := make([]tgbotapi.Chattable, len(chunks))

	for i := 0; i < len(chunks); i++ {
		if i == 0 && m.ImageURL != "" {
			photo := tgbotapi.NewPhoto(m.ChatID, tgbotapi.FileURL(m.ImageURL))
			photo.ParseMode = "HTML"
			photo.Caption = chunks[i]
			messages[i] = photo
			continue
		}

		msg := tgbotapi.NewMessage(m.ChatID, chunks[i])
		msg.DisableWebPagePreview = true
		msg.ParseMode = "HTML"
		messages[i] = msg
	}

	return messages, nil
}
