package task

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/db"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/model"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"golang.org/x/exp/slog"
)

type HandlerTgChat struct {
	logger    *slog.Logger
	publisher common.Pub
	chatRepo  repository.ReadWriteRepository[*entity.Chat]
}

func (h *HandlerTgChat) ProcessTask(ctx context.Context, task *asynq.Task) error {
	switch task.Type() {
	case TelegramChatNew:
		return h.new(ctx, task)
	case TelegramChatEdit:
		return h.edit(ctx, task)
	}
	return nil
}

func (h *HandlerTgChat) new(ctx context.Context, task *asynq.Task) error {
	var chat tgbotapi.Chat
	if err := unmarshal(task.Payload(), &chat); err != nil {
		h.logger.Error("error due to unmarshal task payload", "err", err)
		return nil
	}
	return h.save(ctx, chat, nil)
}

func (h *HandlerTgChat) edit(ctx context.Context, task *asynq.Task) error {
	var chat tgbotapi.ChatMemberUpdated
	if err := unmarshal(task.Payload(), &chat); err != nil {
		h.logger.Error("error due to unmarshal task payload", "err", err)
		return nil
	}
	return h.save(ctx, chat.Chat, &chat.NewChatMember)
}

func (h *HandlerTgChat) save(ctx context.Context, tgChat tgbotapi.Chat, member *tgbotapi.ChatMember) error {
	chat, err := h.toEntityChat(ctx, tgChat)
	if err != nil {
		h.logger.Error("error due to find chat", "err", err, "chat", tgChat, "telegram_id", tgChat.ID)
		return fmt.Errorf("%s %w", OpServerProcessTask, err)
	}

	if member != nil {
		chat.SetDeleted(member.HasLeft() || member.WasKicked())
		chat.SetRights(entity.ChatRights{
			Status:                member.Status,
			IsAnonymous:           member.IsAnonymous,
			UntilDate:             member.UntilDate,
			CanBeEdited:           member.CanBeEdited,
			CanManageChat:         member.CanManageChat,
			CanPostMessages:       member.CanPostMessages,
			CanEditMessages:       member.CanEditMessages,
			CanDeleteMessages:     member.CanDeleteMessages,
			CanRestrictMembers:    member.CanRestrictMembers,
			CanPromoteMembers:     member.CanPromoteMembers,
			CanChangeInfo:         member.CanChangeInfo,
			CanInviteUsers:        member.CanInviteUsers,
			CanPinMessages:        member.CanPinMessages,
			IsMember:              member.IsMember,
			CanSendMessages:       member.CanSendMessages,
			CanSendMediaMessages:  member.CanSendMediaMessages,
			CanSendPolls:          member.CanSendPolls,
			CanSendOtherMessages:  member.CanSendOtherMessages,
			CanAddWebPagePreviews: member.CanAddWebPagePreviews,
		})
	}

	if err = h.chatRepo.Save(ctx, chat); err != nil {
		h.logger.Error("error due to save chat", "err", err, "chat", tgChat, "telegram_id", tgChat.ID)
		return fmt.Errorf("%s %w", OpServerProcessTask, err)
	}

	h.publisher.Telegram(ctx, model.Message{View: model.ViewChat, Data: chat})

	return nil
}

func (h *HandlerTgChat) toEntityChat(ctx context.Context, tgChat tgbotapi.Chat) (*entity.Chat, error) {
	chat, err := h.find(ctx, tgChat.ID)
	if err != nil {
		if err != repository.ErrEntityNotFound {
			return nil, err
		}
		chat = &entity.Chat{
			ID:         uuid.New(),
			TelegramID: tgChat.ID,
		}
		chat.SetBlocked(false)
		chat.SetDeleted(false)
		chat.SetBroadcast([]uuid.UUID{})
	}

	chat.Type = entity.ChatType(tgChat.Type)
	chat.Title = tgChat.Title
	chat.Username = tgChat.UserName
	chat.FirstName = tgChat.FirstName
	chat.LastName = tgChat.LastName

	return chat, nil
}

func (h *HandlerTgChat) find(ctx context.Context, chatID int64) (*entity.Chat, error) {
	criteria := db.BuildCriteria(fmt.Sprintf("size=1&field.0.0=telegram_id&value.0.0=%d", chatID))
	chats, err := h.chatRepo.Find(ctx, criteria)
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return nil, repository.ErrEntityNotFound
	}
	return chats[0], nil
}
