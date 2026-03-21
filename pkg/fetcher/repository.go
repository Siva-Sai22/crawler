package fetcher

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WebsiteRepository struct {
	db *pgxpool.Pool
}

type Website struct {
	ID      int64
	URL     string
	Content string
}

var ErrDuplicate = errors.New("record already exists")

func NewWebsiteRepository(db *pgxpool.Pool) *WebsiteRepository {
	return &WebsiteRepository{
		db: db,
	}
}

func (r *WebsiteRepository) Migrate(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS website (
			id BIGSERIAL PRIMARY KEY,
			url TEXT NOT NULL,
			content TEXT
		);
	`

	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *WebsiteRepository) Create(ctx context.Context, url string, content string) (*Website, error) {
	var id int64
	query := `INSERT INTO website (url, content) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, url, content).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) && pgxError.Code == "23505" {
			return nil, ErrDuplicate
		}
		return nil, err
	}

	return &Website{
		ID:      id,
		URL:     url,
		Content: content,
	}, nil
}
