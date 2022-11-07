package rooms

import "github.com/rumorsflow/rumors/internal/models"

type UpdateRequest struct {
	Broadcast *[]string `json:"broadcast,omitempty" validate:"omitempty,dive,uuid4"`
}

func (r UpdateRequest) Room(id int64) models.Room {
	return models.Room{
		Id:        id,
		Broadcast: r.Broadcast,
	}
}
