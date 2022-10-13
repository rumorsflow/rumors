package models

import "time"

type (
	RoomType       string
	RoomPermission string
)

const (
	Private    RoomType = "private"
	Group      RoomType = "group"
	SuperGroup RoomType = "supergroup"
	Channel    RoomType = "channel"

	BeEdited           RoomPermission = "be_edited"
	ManageChat         RoomPermission = "manage_chat"
	PostMessages       RoomPermission = "post_messages"
	EditMessages       RoomPermission = "edit_messages"
	DeleteMessages     RoomPermission = "delete_messages"
	ManageVoiceChats   RoomPermission = "manage_voice_chats"
	RestrictMembers    RoomPermission = "restrict_members"
	PromoteMembers     RoomPermission = "promote_members"
	ChangeInfo         RoomPermission = "change_info"
	InviteUsers        RoomPermission = "invite_users"
	PinMessages        RoomPermission = "pin_messages"
	SendMessages       RoomPermission = "send_messages"
	SendMediaMessages  RoomPermission = "send_media_messages"
	SendPolls          RoomPermission = "send_polls"
	SendOtherMessages  RoomPermission = "send_other_messages"
	AddWebPagePreviews RoomPermission = "add_web_page_previews"
)

type Room struct {
	Id          int64             `json:"id,omitempty" bson:"_id,omitempty"`
	Type        RoomType          `json:"type,omitempty" bson:"type,omitempty"`
	Title       string            `json:"title,omitempty" bson:"title,omitempty"`
	UserName    string            `json:"username,omitempty" bson:"username,omitempty"`
	FirstName   string            `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName    string            `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Broadcast   *[]string         `json:"broadcast,omitempty" bson:"broadcast,omitempty"`
	Permissions *[]RoomPermission `json:"permissions,omitempty" bson:"permissions,omitempty"`
	Deleted     *bool             `json:"deleted,omitempty" bson:"deleted,omitempty"`
	CreatedAt   time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (r *Room) SetBroadcast(broadcast []string) *Room {
	r.Broadcast = &broadcast
	return r
}

func (r *Room) SetPermissions(permissions []RoomPermission) *Room {
	r.Permissions = &permissions
	return r
}

func (r *Room) SetDeleted(deleted bool) *Room {
	r.Deleted = &deleted
	return r
}

func (r *Room) IsDeleted() bool {
	return r.Deleted == nil || *r.Deleted
}
