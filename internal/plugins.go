package internal

import (
	"github.com/rumorsflow/rumors/v2/internal/db"
	"github.com/rumorsflow/rumors/v2/internal/http"
	"github.com/rumorsflow/rumors/v2/internal/pubsub"
	"github.com/rumorsflow/rumors/v2/internal/rdb"
	"github.com/rumorsflow/rumors/v2/internal/task"
	"github.com/rumorsflow/rumors/v2/internal/telegram"
)

func Plugins() []any {
	return []any{
		&rdb.Plugin{},
		&db.Plugin{},
		&pubsub.Plugin{},
		&telegram.Plugin{},
		&task.Plugin{},
		&http.Plugin{},
	}
}
