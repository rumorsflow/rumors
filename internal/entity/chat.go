package entity

import (
	"github.com/google/uuid"
	"time"
)

type ChatType string

const (
	Private    ChatType = "private"
	Group      ChatType = "group"
	SuperGroup ChatType = "supergroup"
	Channel    ChatType = "channel"
)

type ChatRights struct {
	Status                string `json:"status,omitempty" bson:"status,omitempty"`
	IsAnonymous           bool   `json:"is_anonymous,omitempty" bson:"is_anonymous,omitempty"`
	UntilDate             int64  `json:"until_date,omitempty" bson:"until_date,omitempty"`
	CanBeEdited           bool   `json:"can_be_edited,omitempty" bson:"can_be_edited,omitempty"`
	CanManageChat         bool   `json:"can_manage_chat,omitempty" bson:"can_manage_chat,omitempty"`
	CanPostMessages       bool   `json:"can_post_messages,omitempty" bson:"can_post_messages,omitempty"`
	CanEditMessages       bool   `json:"can_edit_messages,omitempty" bson:"can_edit_messages,omitempty"`
	CanDeleteMessages     bool   `json:"can_delete_messages,omitempty" bson:"can_delete_messages,omitempty"`
	CanRestrictMembers    bool   `json:"can_restrict_members,omitempty" bson:"can_restrict_members,omitempty"`
	CanPromoteMembers     bool   `json:"can_promote_members,omitempty" bson:"can_promote_members,omitempty"`
	CanChangeInfo         bool   `json:"can_change_info,omitempty" bson:"can_change_info,omitempty"`
	CanInviteUsers        bool   `json:"can_invite_users,omitempty" bson:"can_invite_users,omitempty"`
	CanPinMessages        bool   `json:"can_pin_messages,omitempty" bson:"can_pin_messages,omitempty"`
	IsMember              bool   `json:"is_member,omitempty" bson:"is_member,omitempty"`
	CanSendMessages       bool   `json:"can_send_messages,omitempty" bson:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool   `json:"can_send_media_messages,omitempty" bson:"can_send_media_messages,omitempty"`
	CanSendPolls          bool   `json:"can_send_polls,omitempty" bson:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool   `json:"can_send_other_messages,omitempty" bson:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews,omitempty" bson:"can_add_web_page_previews,omitempty"`
}

func (chat ChatRights) IsCreator() bool { return chat.Status == "creator" }

func (chat ChatRights) IsAdministrator() bool { return chat.Status == "administrator" }

func (chat ChatRights) HasLeft() bool { return chat.Status == "left" }

func (chat ChatRights) WasKicked() bool { return chat.Status == "kicked" }

type Chat struct {
	ID         uuid.UUID    `json:"id,omitempty" bson:"_id,omitempty"`
	TelegramID int64        `json:"telegram_id,omitempty" bson:"telegram_id,omitempty"`
	Type       ChatType     `json:"type,omitempty" bson:"type,omitempty"`
	Title      string       `json:"title,omitempty" bson:"title,omitempty"`
	Username   string       `json:"username,omitempty" bson:"username,omitempty"`
	FirstName  string       `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName   string       `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Broadcast  *[]uuid.UUID `json:"broadcast,omitempty" bson:"broadcast,omitempty"`
	Rights     *ChatRights  `json:"rights,omitempty" bson:"rights,omitempty"`
	Blocked    *bool        `json:"blocked,omitempty" bson:"blocked,omitempty"`
	Deleted    *bool        `json:"deleted,omitempty" bson:"deleted,omitempty"`
	CreatedAt  time.Time    `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (e *Chat) EntityID() uuid.UUID {
	return e.ID
}

func (e *Chat) SetBroadcast(broadcast []uuid.UUID) *Chat {
	e.Broadcast = &broadcast
	return e
}

func (e *Chat) SetRights(rights ChatRights) *Chat {
	e.Rights = &rights
	return e
}

func (e *Chat) SetBlocked(blocked bool) *Chat {
	e.Blocked = &blocked
	return e
}

func (e *Chat) IsBlocked() bool {
	return e.Blocked != nil && *e.Blocked
}

func (e *Chat) SetDeleted(deleted bool) *Chat {
	e.Deleted = &deleted
	return e
}

func (e *Chat) IsDeleted() bool {
	return e.Deleted != nil && *e.Deleted
}
