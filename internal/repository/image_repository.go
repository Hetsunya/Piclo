package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"Piclo/internal/model"
)

type ImageRepository struct {
	pool *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{pool: pool}
}

func (r *ImageRepository) Save(ctx context.Context, img *model.Image) error {
	query := `
        INSERT INTO images (id, filename, mime_type, size, storage_key, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query,
		img.ID, img.Filename, img.MIMEType, img.Size, img.StorageKey, img.CreatedAt, img.ExpiresAt)
	return err
}

func (r *ImageRepository) GetByID(ctx context.Context, id string) (*model.Image, error) {
	img := &model.Image{}
	query := `SELECT id, filename, mime_type, size, storage_key, created_at, expires_at 
              FROM images WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&img.ID, &img.Filename, &img.MIMEType, &img.Size,
		&img.StorageKey, &img.CreatedAt, &img.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return img, nil
}
