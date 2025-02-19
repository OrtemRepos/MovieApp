package postgresql

import (
	"context"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

var schema = `
CREATE TABLE IF NOS EXISTS ratings (
	record_id VARCHAR(255),
	record_type VARCHAR(255),
	user_id VARCHAR(255),
	value INT
);
`

type Repository struct {
	db *sqlx.DB
}

func New() (*Repository, error) {
	db, err := sqlx.Open("pgx", "user=root dbname=movie sslmode=disable")
	if err != nil {
		return nil, err
	}
	db.MustExec(schema)
	return &Repository{db}, nil
}

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	rows, err := r.db.QueryxContext(
		ctx, "SELECT user_id, value FROM ratings WHERE record_id = ? AND record_type = ?",
		recordID, recordType,
	)
	if err != nil {
		return nil, err
	}
	var res []model.Rating
	for rows.Next() {
		var userID string
		var value int32
		if err := rows.Scan(&userID, &value); err != nil {
			return nil, err
		}
		res = append(
			res,
			model.Rating{
				RecordID: recordID,
				RecordType: recordType,
				UserID: model.UserID(userID),
				Value: model.RatingValue(value),
			},
		)
	}
	if len(res) == 0 {
		return nil, repository.ErrNotFound
	}
	return res, nil
}

func (r *Repository) Put(
	ctx context.Context, recordID model.RecordID,
	recordType model.RecordType, rating *model.Rating,
) error {
	_, err := r.db.ExecContext(
		ctx, "INSERT INTO ratings (record_id, record_type, user_id, value) VALUES (?, ?, ?, ?)",
		recordID, recordType, rating.UserID, rating.Value,
	)
	return err
}