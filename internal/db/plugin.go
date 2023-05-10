package db

import (
	"context"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/pkg/config"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"reflect"
	"sync"
)

const PluginName = "mongo"

type Plugin struct {
	resolvers sync.Map
}

func (p *Plugin) Init(cfg config.Configurer) error {
	const op = errors.Op("db_plugin_init")

	if !cfg.Has(PluginName) {
		return errors.E(op, errors.Disabled)
	}

	var c mongodb.Config
	if err := cfg.UnmarshalKey(PluginName, &c); err != nil {
		return errors.E(op, err)
	}

	database, err := mongodb.NewDatabase(context.Background(), &c)
	if err != nil {
		return errors.E(op, err)
	}

	p.resolvers.Store((*entity.Site)(nil), newResolver[*entity.Site](func() (repository.ReadWriteRepository[*entity.Site], error) {
		return NewRepository[*entity.Site](
			database,
			entity.SiteCollection,
			WithEntityFactory(repository.Factory[*entity.Site]()),
			WithBeforeSave(BeforeSave[*entity.Site]),
			WithAfterSave(AfterSave[*entity.Site]),
			WithIndexes[*entity.Site](SiteIndexes),
		)
	}))

	p.resolvers.Store((*entity.Article)(nil), newResolver[*entity.Article](func() (repository.ReadWriteRepository[*entity.Article], error) {
		return NewRepository[*entity.Article](
			database,
			entity.ArticleCollection,
			WithEntityFactory(repository.Factory[*entity.Article]()),
			WithBeforeSave(ArticleBeforeSave),
			WithAfterSave(AfterSave[*entity.Article]),
			WithIndexes[*entity.Article](ArticleIndexes),
		)
	}))

	p.resolvers.Store((*entity.Chat)(nil), newResolver[*entity.Chat](func() (repository.ReadWriteRepository[*entity.Chat], error) {
		return NewRepository[*entity.Chat](
			database,
			entity.ChatCollection,
			WithEntityFactory(repository.Factory[*entity.Chat]()),
			WithBeforeSave(ChatBeforeSave),
			WithAfterSave(AfterSave[*entity.Chat]),
			WithIndexes[*entity.Chat](ChatIndexes),
		)
	}))

	p.resolvers.Store((*entity.Job)(nil), newResolver[*entity.Job](func() (repository.ReadWriteRepository[*entity.Job], error) {
		return NewRepository[*entity.Job](
			database,
			entity.JobCollection,
			WithEntityFactory(repository.Factory[*entity.Job]()),
			WithBeforeSave(BeforeSave[*entity.Job]),
			WithAfterSave(AfterSave[*entity.Job]),
			WithIndexes[*entity.Job](JobIndexes),
		)
	}))

	p.resolvers.Store((*entity.SysUser)(nil), newResolver[*entity.SysUser](func() (repository.ReadWriteRepository[*entity.SysUser], error) {
		return NewRepository[*entity.SysUser](
			database,
			entity.SysUserCollection,
			WithEntityFactory(repository.Factory[*entity.SysUser]()),
			WithBeforeSave(BeforeSave[*entity.SysUser]),
			WithAfterSave(AfterSave[*entity.SysUser]),
			WithIndexes[*entity.SysUser](SysUserIndexes),
		)
	}))

	return nil
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*common.UnitOfWork)(nil), p.ServiceUnitOfWork),
	}
}

func (p *Plugin) ServiceUnitOfWork() common.UnitOfWork {
	return p
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Repository(tp any) (any, error) {
	const op = errors.Op("repository_resolver")

	if r, ok := p.resolvers.Load(tp); ok {
		resp := reflect.ValueOf(r).MethodByName("Resolve").Call([]reflect.Value{})
		if err, ok := resp[1].Interface().(error); ok {
			return nil, errors.E(op, err)
		}
		return resp[0].Interface(), nil
	}

	return nil, errors.E(op, errors.Errorf("repository.ReadWriteRepository[%T] not found", tp))
}
