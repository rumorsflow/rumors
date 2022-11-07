package rooms

import (
	"github.com/alexedwards/flow"
	"github.com/rumorsflow/rumors/internal/api/middleware/authz"
	"github.com/rumorsflow/rumors/internal/api/util"
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/spf13/cast"
	"net/http"
)

const (
	PluginName = "/api/v1/rooms"
	id         = "/:id"
)

type Plugin struct {
	storage storage.RoomStorage
}

func (p *Plugin) Init(s storage.RoomStorage) error {
	p.storage = s
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Register(mux *flow.Mux) {
	mux.HandleFunc(PluginName, authz.IsAdmin(p.list), http.MethodGet)
	mux.HandleFunc(PluginName+id, authz.IsAdmin(p.read), http.MethodGet)
	mux.HandleFunc(PluginName+id, authz.IsAdmin(p.update), http.MethodPatch)
}

func (p *Plugin) list(w http.ResponseWriter, r *http.Request) {
	criteria := util.GetCriteria(r)

	data, err := p.storage.Find(r.Context(), criteria)
	if err != nil {
		panic(err)
	}

	total, err := p.storage.Count(r.Context(), criteria.Filter)
	if err != nil {
		panic(err)
	}

	util.OK(w, util.ListResponse{
		Data:  data,
		Total: total,
		Index: criteria.Index,
		Size:  criteria.Size,
	})
}

func (p *Plugin) read(w http.ResponseWriter, r *http.Request) {
	data, err := p.storage.FindById(r.Context(), getId(r))
	if err != nil {
		panic(err)
	}
	util.OK(w, data)
}

func (p *Plugin) update(w http.ResponseWriter, r *http.Request) {
	var dto UpdateRequest
	util.Bind(r, &dto)

	model := dto.Room(getId(r))
	if err := p.storage.Save(r.Context(), &model); err != nil {
		panic(err)
	}
	util.NoContent(w)
}

func getId(r *http.Request) int64 {
	return cast.ToInt64(util.GetId(r))
}
