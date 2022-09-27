package notifications

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type tgbot struct {
	log   *zerolog.Logger
	bot   *tgbotapi.BotAPI
	owner int64
}

func NewTgBotNotification(owner int64, bot *tgbotapi.BotAPI, log *zerolog.Logger) Notification {
	return &tgbot{
		log:   log,
		bot:   bot,
		owner: owner,
	}
}

func (n *tgbot) Forbidden(to any) {
	n.Send(to, "<b>Forbidden</b>\n\nExecute access forbidden")
}

func (n *tgbot) Success(to any, text string) {
	n.log.Info().Interface("to", to).Msg(text)
	n.Send(to, "<b>Success</b>\n\n"+text)
}

func (n *tgbot) Error(to any, text string) {
	n.log.Error().Interface("to", to).Str("error", text).Msg("")
	n.Send(to, "<b>Error</b>\n\n"+text)
}

func (n *tgbot) Err(to any, err error) {
	n.Error(to, err.Error())
}

func (n *tgbot) Send(to any, text string, options ...any) {
	chatID := n.owner
	if id, ok := to.(int64); ok && id != 0 {
		chatID = id
	}

	disableWebPagePreview := true
	parseMode := "html"

	for _, opt := range options {
		switch o := opt.(type) {
		case bool:
			disableWebPagePreview = o
		case string:
			parseMode = o
		}
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableWebPagePreview = disableWebPagePreview
	msg.ParseMode = parseMode

	n.raw(msg)
}

func (n *tgbot) raw(c tgbotapi.Chattable) {
	if _, err := n.bot.Send(c); err != nil {
		n.log.Error().Err(err).Msg("")
	}
}
