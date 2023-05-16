package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"github.com/rumorsflow/rumors/v2/pkg/util"
)

const OpMessageToChattableList = "message: to chattable list ->"

func chattableList(m model.Message, render model.Render, ownerID int64) ([]tgbotapi.Chattable, error) {
	chatID := ownerID
	if m.ChatID != 0 {
		chatID = m.ChatID
	}

	if m.View == "" {
		return nil, fmt.Errorf("%s error: message view is required", OpMessageToChattableList)
	}

	text, err := render(m.View, m.Data)
	if err != nil {
		return nil, fmt.Errorf("%s execute template error: %w", OpMessageToChattableList, err)
	}

	chunks := util.SplitMax(text, "\n", 4096)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("%s error: split text in chunks", OpMessageToChattableList)
	}

	messages := make([]tgbotapi.Chattable, len(chunks))

	for i := 0; i < len(chunks); i++ {
		if i == 0 && m.ImageURL != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(m.ImageURL))
			photo.ParseMode = "HTML"
			photo.Caption = chunks[i]
			messages[i] = photo
			continue
		}

		msg := tgbotapi.NewMessage(chatID, chunks[i])
		msg.DisableWebPagePreview = true
		msg.ParseMode = "HTML"
		messages[i] = msg
	}

	return messages, nil
}
