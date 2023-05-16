package model

import (
	"github.com/goccy/go-json"
	"github.com/rumorsflow/rumors/v2/internal/entity"
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
)

type Render func(view View, data any) (string, error)

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
		return m.unmarshalData(msg.Data, map[string][]Article{})
	case ViewChat:
		return m.unmarshalData(msg.Data, entity.Chat{})
	case ViewSites, ViewSub:
		return m.unmarshalData(msg.Data, []string{})
	default:
		return m.unmarshalData(msg.Data, "")
	}
}

func (m *Message) unmarshalData(data json.RawMessage, i any) error {
	if err := json.Unmarshal(data, i); err != nil {
		return err
	}
	m.Data = i
	return nil
}
