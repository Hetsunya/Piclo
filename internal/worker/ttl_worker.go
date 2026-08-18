package worker

import (
	"context"
	"log"
	"time"

	"Piclo/internal/repository"
	"Piclo/internal/storage"
)

type TTLWorker struct {
	repo     *repository.ImageRepository
	store    storage.Storage
	interval time.Duration
}

func NewTTLWorker(repo *repository.ImageRepository, store storage.Storage, interval time.Duration) *TTLWorker {
	return &TTLWorker{
		repo:     repo,
		store:    store,
		interval: interval,
	}
}

// Run starts background cleanup. Blocks execution until context cancellation.
func (w *TTLWorker) Run(ctx context.Context) {
	log.Printf("TTL Worker started with interval: %v", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("TTL Worker stopped gracefully")
			return
		case <-ticker.C:
			w.cleanup(ctx)
		}
	}
}

func (w *TTLWorker) cleanup(ctx context.Context) {
	// 1. Find expired records
	expiredImages, err := w.repo.GetExpired(ctx)
	if err != nil {
		log.Printf("TTL Worker: failed to get expired images: %v", err)
		return
	}

	if len(expiredImages) == 0 {
		return // Nothing to delete
	}

	log.Printf("TTL Worker: found %d expired images, starting cleanup...", len(expiredImages))

	// 2. Delete each one
	for _, img := range expiredImages {
		// First delete from MinIO.
		// IMPORTANT: If MinIO fails, the DB record will remain, and worker will retry in a minute (idempotent).
		if err := w.store.Delete(ctx, img.StorageKey); err != nil {
			log.Printf("TTL Worker: failed to delete from storage %s: %v", img.StorageKey, err)
			continue // Skip DB deletion, try again later
		}

		// Then delete from PostgreSQL
		if err := w.repo.Delete(ctx, img.ID); err != nil {
			log.Printf("TTL Worker: failed to delete from DB id %s: %v", img.ID, err)
		} else {
			log.Printf("TTL Worker: successfully deleted image %s", img.ID)
		}
	}
}
