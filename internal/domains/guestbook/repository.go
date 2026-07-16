package guestbook

import (
	"context"
	"ekel-backend/internal/models"
)

type GuestbookRepository interface {
	Create(ctx context.Context, entry *models.GuestbookEntry) error
	GetByID(ctx context.Context, id string) (*models.GuestbookEntry, error)
	List(ctx context.Context) ([]*models.GuestbookEntry, error)
	Update(ctx context.Context, entry *models.GuestbookEntry) error
	Delete(ctx context.Context, id string) error
}
