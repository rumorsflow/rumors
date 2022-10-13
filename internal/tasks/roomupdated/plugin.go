package roomupdated

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/rumorsflow/rumors/internal/consts"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/services/room"
	"go.uber.org/zap"
)

const (
	PluginName = consts.TaskRoomUpdated

	statusLeft   = "left"
	statusKicked = "kicked"
)

type Plugin struct {
	log     *zap.Logger
	service room.Service
}

func (p *Plugin) Init(log *zap.Logger, service room.Service) error {
	p.log = log
	p.service = service
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var member tgbotapi.ChatMemberUpdated
	if err := json.Unmarshal(task.Payload(), &member); err != nil {
		p.log.Error("error due to unmarshal task payload", zap.Error(err))
		return nil
	}

	deleted := member.NewChatMember.Status == statusLeft || member.NewChatMember.Status == statusKicked

	return p.service.ChatMemberUpdated(ctx, member.Chat, deleted, toPermissions(member.NewChatMember)...)
}

func toPermissions(member tgbotapi.ChatMember) []models.RoomPermission {
	var permissions []models.RoomPermission
	if member.CanBeEdited {
		permissions = append(permissions, models.BeEdited)
	}
	if member.CanManageChat {
		permissions = append(permissions, models.ManageChat)
	}
	if member.CanPostMessages {
		permissions = append(permissions, models.PostMessages)
	}
	if member.CanEditMessages {
		permissions = append(permissions, models.EditMessages)
	}
	if member.CanDeleteMessages {
		permissions = append(permissions, models.DeleteMessages)
	}
	if member.CanManageVoiceChats {
		permissions = append(permissions, models.ManageVoiceChats)
	}
	if member.CanRestrictMembers {
		permissions = append(permissions, models.RestrictMembers)
	}
	if member.CanPromoteMembers {
		permissions = append(permissions, models.PromoteMembers)
	}
	if member.CanChangeInfo {
		permissions = append(permissions, models.ChangeInfo)
	}
	if member.CanInviteUsers {
		permissions = append(permissions, models.InviteUsers)
	}
	if member.CanPinMessages {
		permissions = append(permissions, models.PinMessages)
	}
	if member.CanSendMessages {
		permissions = append(permissions, models.SendMessages)
	}
	if member.CanSendMediaMessages {
		permissions = append(permissions, models.SendMediaMessages)
	}
	if member.CanSendPolls {
		permissions = append(permissions, models.SendPolls)
	}
	if member.CanSendOtherMessages {
		permissions = append(permissions, models.SendOtherMessages)
	}
	if member.CanAddWebPagePreviews {
		permissions = append(permissions, models.AddWebPagePreviews)
	}
	return permissions
}
