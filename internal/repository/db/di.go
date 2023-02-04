package db

import (
	"context"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/errs"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
)

type (
	ArticleKey struct{}
	ChatKey    struct{}
	FeedKey    struct{}
	JobKey     struct{}
	SysUserKey struct{}
)

func NewArticleRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Article], error) {
	return di.New[*Repository[*entity.Article]](ctx, ArticleKey{}, c...)
}

func NewChatRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Chat], error) {
	return di.New[*Repository[*entity.Chat]](ctx, ChatKey{}, c...)
}

func NewFeedRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Feed], error) {
	return di.New[*Repository[*entity.Feed]](ctx, FeedKey{}, c...)
}

func NewJobRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Job], error) {
	return di.New[*Repository[*entity.Job]](ctx, JobKey{}, c...)
}

func NewSysUserRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.SysUser], error) {
	return di.New[*Repository[*entity.SysUser]](ctx, SysUserKey{}, c...)
}

func GetArticleRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Article], error) {
	return di.Get[*Repository[*entity.Article]](ctx, ArticleKey{}, c...)
}

func GetChatRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Chat], error) {
	return di.Get[*Repository[*entity.Chat]](ctx, ChatKey{}, c...)
}

func GetFeedRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Feed], error) {
	return di.Get[*Repository[*entity.Feed]](ctx, FeedKey{}, c...)
}

func GetJobRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.Job], error) {
	return di.Get[*Repository[*entity.Job]](ctx, JobKey{}, c...)
}

func GetSysUserRepository(ctx context.Context, c ...di.Container) (*Repository[*entity.SysUser], error) {
	return di.Get[*Repository[*entity.SysUser]](ctx, SysUserKey{}, c...)
}

func ArticleActivator() *di.Activator {
	return &di.Activator{
		Key: ArticleKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			database, err := mongodb.Get(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return ToNilCloser(NewRepository[*entity.Article](
				database,
				"articles",
				WithEntityFactory(repository.Factory[*entity.Article]()),
				WithBeforeSave(ArticleBeforeSave),
				WithAfterSave(AfterSave[*entity.Article]),
				WithIndexes[*entity.Article](ArticleIndexes),
			))
		}),
	}
}

func ChatActivator() *di.Activator {
	return &di.Activator{
		Key: ChatKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			database, err := mongodb.Get(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return ToNilCloser(NewRepository[*entity.Chat](
				database,
				"chats",
				WithEntityFactory(repository.Factory[*entity.Chat]()),
				WithBeforeSave(ChatBeforeSave),
				WithAfterSave(AfterSave[*entity.Chat]),
				WithIndexes[*entity.Chat](ChatIndexes),
			))
		}),
	}
}

func FeedActivator() *di.Activator {
	return &di.Activator{
		Key: FeedKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			database, err := mongodb.Get(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return ToNilCloser(NewRepository[*entity.Feed](
				database,
				"feeds",
				WithEntityFactory(repository.Factory[*entity.Feed]()),
				WithBeforeSave(BeforeSave[*entity.Feed]),
				WithAfterSave(AfterSave[*entity.Feed]),
				WithIndexes[*entity.Feed](FeedIndexes),
			))
		}),
	}
}

func JobActivator() *di.Activator {
	return &di.Activator{
		Key: JobKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			database, err := mongodb.Get(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return ToNilCloser(NewRepository[*entity.Job](
				database,
				"jobs",
				WithEntityFactory(repository.Factory[*entity.Job]()),
				WithBeforeSave(BeforeSave[*entity.Job]),
				WithAfterSave(AfterSave[*entity.Job]),
				WithIndexes[*entity.Job](JobIndexes),
			))
		}),
	}
}

func SysUserActivator() *di.Activator {
	return &di.Activator{
		Key: SysUserKey{},
		Factory: di.FactoryFunc(func(ctx context.Context, c di.Container) (any, di.Closer, error) {
			database, err := mongodb.Get(ctx, c)
			if err != nil {
				return nil, nil, err
			}

			return ToNilCloser(NewRepository[*entity.SysUser](
				database,
				"sys_users",
				WithEntityFactory(repository.Factory[*entity.SysUser]()),
				WithBeforeSave(BeforeSave[*entity.SysUser]),
				WithAfterSave(AfterSave[*entity.SysUser]),
				WithIndexes[*entity.SysUser](SysUserIndexes),
			))
		}),
	}
}

func ToNilCloser[T any](value T, err error) (T, di.Closer, error) {
	if err != nil {
		err = errs.E(di.OpFactory, err)
	}
	return value, nil, err
}
