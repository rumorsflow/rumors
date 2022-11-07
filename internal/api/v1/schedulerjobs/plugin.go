package schedulerjobs

import (
	"fmt"
	"github.com/alexedwards/flow"
	"github.com/rumorsflow/rumors/internal/api/middleware/authz"
	"github.com/rumorsflow/rumors/internal/api/util"
	"github.com/rumorsflow/scheduler-mongo-provider"
	"net/http"
)

const (
	PluginName = "/api/v1/jobs"
	id         = "/:id"
)

type Plugin struct {
	storage smp.PeriodicTaskStorage
}

func (p *Plugin) Init(s smp.PeriodicTaskStorage) error {
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
	mux.HandleFunc(PluginName, authz.IsAdmin(p.create), http.MethodPost)
	mux.HandleFunc(PluginName+id, authz.IsAdmin(p.update), http.MethodPatch)
	mux.HandleFunc(PluginName+id, authz.IsAdmin(p.delete), http.MethodDelete)
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
	data, err := p.storage.FindById(r.Context(), util.GetId(r))
	if err != nil {
		panic(err)
	}
	util.OK(w, data)
}

func (p *Plugin) create(w http.ResponseWriter, r *http.Request) {
	var dto CreateRequest
	util.Bind(r, &dto)

	model := dto.PeriodicTask()
	if err := p.storage.Save(r.Context(), &model); err != nil {
		panic(err)
	}

	util.Created(w, fmt.Sprintf("%s/%s", PluginName, model.Id))
}

func (p *Plugin) update(w http.ResponseWriter, r *http.Request) {
	var dto UpdateRequest
	util.Bind(r, &dto)

	model := dto.PeriodicTask(util.GetId(r))
	if err := p.storage.Save(r.Context(), &model); err != nil {
		panic(err)
	}
	util.NoContent(w)
}

func (p *Plugin) delete(w http.ResponseWriter, r *http.Request) {
	if err := p.storage.Delete(r.Context(), util.GetId(r)); err != nil {
		panic(err)
	}
	util.NoContent(w)
}
