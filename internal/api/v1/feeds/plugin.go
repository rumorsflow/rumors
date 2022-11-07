package feeds

import (
	"fmt"
	"github.com/alexedwards/flow"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/api/util"
	"github.com/rumorsflow/rumors/internal/storage"
	"github.com/rumorsflow/rumors/internal/tgbotsender"
	"net/http"
)

const (
	PluginName = "/api/v1/feeds"
	id         = "/:id"
)

type Plugin struct {
	storage storage.FeedStorage
	sender  tgbotsender.TelegramSender
}

func (p *Plugin) Init(s storage.FeedStorage, sender tgbotsender.TelegramSender) error {
	p.storage = s
	p.sender = sender
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Register(mux *flow.Mux) {
	mux.HandleFunc(PluginName, p.list, http.MethodGet)
	mux.HandleFunc(PluginName+id, p.read, http.MethodGet)
	mux.HandleFunc(PluginName, p.create, http.MethodPost)
	mux.HandleFunc(PluginName+id, p.update, http.MethodPatch)
	mux.HandleFunc(PluginName+id, p.delete, http.MethodDelete)
}

func (p *Plugin) list(w http.ResponseWriter, r *http.Request) {
	criteria, _ := mongoext.GetC(r.URL.RawQuery, "filters")

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

	model := dto.Feed(p.sender.Owner())
	if err := p.storage.Save(r.Context(), &model); err != nil {
		panic(err)
	}

	util.Created(w, fmt.Sprintf("%s/%s", PluginName, model.Id))
}

func (p *Plugin) update(w http.ResponseWriter, r *http.Request) {
	var dto UpdateRequest
	util.Bind(r, &dto)

	model := dto.Feed(util.GetId(r))
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
