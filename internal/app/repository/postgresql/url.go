package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"simple-url-shortener/internal/app/db"
	"simple-url-shortener/internal/app/repository"
)

type UrlRepo struct {
	db db.DBops
}

func (u UrlRepo) Add(ctx context.Context, Url *repository.Url) (int64, error) {
	var id int64
	err := u.db.ExecQueryRow(ctx, `INSERT INTO urls(alias,url) VALUES($1,$2) RETURNING id;`, Url.Alias, Url.URL).Scan(&id)
	return id, err
}

func (u UrlRepo) GetByID(ctx context.Context, id int64) (*repository.Url, error) {
	var a repository.Url
	err := u.db.Get(ctx, &a, "SELECT id,alias,url FROM urls WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (u UrlRepo) GetByAlias(ctx context.Context, alias string) (*repository.Url, error) {
	var a repository.Url
	err := u.db.Get(ctx, &a, "SELECT id,alias,url FROM urls WHERE alias=$1", alias)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (u UrlRepo) GetAliasByURL(ctx context.Context, urlToFind string) (string, error) {
	var a repository.Url
	err := u.db.Get(ctx, &a, "SELECT id,alias,url FROM urls WHERE url=$1", urlToFind)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrObjectNotFound
		}
		return "", err
	}
	return a.Alias, nil
}

func (u UrlRepo) GetAll(ctx context.Context) (*[]repository.Url, error) {
	var urls []repository.Url
	err := u.db.Select(ctx, &urls, "SELECT id, alias, url FROM urls")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &urls, nil
}

func (u UrlRepo) CheckAliasExists(ctx context.Context, alias string) (bool, error) {
	var count int
	err := u.db.Get(ctx, &count, "SELECT COUNT(*) FROM urls WHERE alias = $1", alias)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u UrlRepo) CheckURLExists(ctx context.Context, urlToCheck string) (bool, error) {
	var count int
	err := u.db.Get(ctx, &count, "SELECT COUNT(*) FROM urls WHERE url = $1", urlToCheck)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u UrlRepo) DeleteByID(ctx context.Context, id int64) error {
	result, err := u.db.Exec(ctx, "DELETE FROM urls WHERE id=$1", id)
	rowsAffected := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrObjectNotFound
	}
	return err
}

func (u UrlRepo) DeleteByAlias(ctx context.Context, alias string) error {
	_, err := u.db.Exec(ctx, "DELETE FROM urls WHERE alias = $1", alias)
	if err != nil {
		return err
	}
	return nil
}

func (u UrlRepo) Update(ctx context.Context, url *repository.Url) error {
	result, err := u.db.Exec(ctx, "UPDATE urls SET alias=$1, url=$2 WHERE id=$3", url.Alias, url.URL, url.ID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return repository.ErrObjectNotFound
	}
	return nil
}

func NewUrls(database db.DBops) *UrlRepo {
	return &UrlRepo{db: database}
}
