//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
package repository

import "context"

type URLsRepo interface {
	Add(ctx context.Context, Url *Url) (int64, error)
	GetByID(ctx context.Context, id int64) (*Url, error)
	GetByAlias(ctx context.Context, alias string) (*Url, error)
	GetAliasByURL(ctx context.Context, urlToFind string) (string, error)
	GetAll(ctx context.Context) (*[]Url, error)
	CheckAliasExists(ctx context.Context, alias string) (bool, error)
	CheckURLExists(ctx context.Context, urlToCheck string) (bool, error)
	DeleteByID(ctx context.Context, id int64) error
	DeleteByAlias(ctx context.Context, alias string) error
	Update(ctx context.Context, url *Url) error
}
