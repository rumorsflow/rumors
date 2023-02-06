package sys

import (
	"github.com/google/uuid"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/http/action"
	"github.com/rumorsflow/rumors/v2/internal/repository"
)

type CreateChatDTO struct {
	TelegramID int64           `json:"telegram_id,omitempty" validate:"required"`
	Type       entity.ChatType `json:"type,omitempty" validate:"required,oneof=private group supergroup channel"`
	Title      string          `json:"title,omitempty" validate:"omitempty,max=254"`
	UserName   string          `json:"user_name,omitempty" validate:"omitempty,max=254"`
	FirstName  string          `json:"first_name,omitempty" validate:"omitempty,max=254"`
	LastName   string          `json:"last_name,omitempty" validate:"omitempty,max=254"`
	Broadcast  []string        `json:"broadcast,omitempty" validate:"omitempty,dive,uuid4"`
	Blocked    bool            `json:"blocked,omitempty"`
	Deleted    bool            `json:"deleted,omitempty"`
}

func (dto CreateChatDTO) toEntity(id uuid.UUID) *entity.Chat {
	broadcast := make([]uuid.UUID, len(dto.Broadcast))
	for i, b := range dto.Broadcast {
		broadcast[i] = uuid.MustParse(b)
	}

	return &entity.Chat{
		ID:         id,
		TelegramID: dto.TelegramID,
		Type:       dto.Type,
		Title:      dto.Title,
		UserName:   dto.UserName,
		FirstName:  dto.FirstName,
		LastName:   dto.LastName,
		Broadcast:  &broadcast,
		Blocked:    &dto.Blocked,
		Deleted:    &dto.Deleted,
	}
}

type UpdateChatDTO struct {
	Broadcast *[]string `json:"broadcast,omitempty" validate:"omitempty,dive,uuid4"`
	Blocked   *bool     `json:"blocked,omitempty"`
}

func (dto UpdateChatDTO) toEntity(id uuid.UUID) *entity.Chat {
	m := &entity.Chat{
		ID: id,
	}
	if dto.Broadcast != nil {
		broadcast := make([]uuid.UUID, len(*dto.Broadcast))
		for i, b := range *dto.Broadcast {
			broadcast[i] = uuid.MustParse(b)
		}
		m.SetBroadcast(broadcast)
	}
	if m.Blocked != nil {
		m.SetBlocked(*m.Blocked)
	}
	return m
}

func NewChatCRUD(
	read repository.ReadRepository[*entity.Chat],
	write repository.WriteRepository[*entity.Chat],
) action.CRUD {
	return action.NewCRUD[*CreateChatDTO, *UpdateChatDTO, *entity.Chat, any](
		read,
		write,
		action.NewDTOFactory[*CreateChatDTO](),
		action.NewDTOFactory[*UpdateChatDTO](),
		action.RequestMapperFunc[*CreateChatDTO, *entity.Chat](func(id uuid.UUID, dto *CreateChatDTO) (*entity.Chat, error) {
			return dto.toEntity(id), nil
		}),
		action.RequestMapperFunc[*UpdateChatDTO, *entity.Chat](func(id uuid.UUID, dto *UpdateChatDTO) (*entity.Chat, error) {
			return dto.toEntity(id), nil
		}),
		nil,
	)
}
