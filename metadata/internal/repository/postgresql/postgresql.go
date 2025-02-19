package postgresql

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

type Repository struct {
	db *sqlx.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS movies (
	id VARCHAR(255),
	title VARCHAR(255),
	description TEXT,
	director VARCHAR(255)
);
`

func New() (*Repository, error) {
	db, err := sqlx.Open("pgx", "user=root dbname=movie sslmode=disable")
	if err != nil {
		return nil, err
	}
	db.MustExec(schema)
	return &Repository{db}, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowxContext(ctx, "SELECT title, description, director FROM movies WHERE id = ?", id)
	if err := row.Scan(&title, description, director); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	m := model.Metadata{ID: id, Title: title, Description: description, Director: director}
	return &m, nil
}

func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)",
		id, metadata.Title, metadata.Description, metadata.Director, 
	)
	return err
}