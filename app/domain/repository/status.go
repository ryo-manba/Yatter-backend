package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Find Status
	FindWithAccountByID(ctx context.Context, id object.StatusID) (*object.Status, error)
	// Create Status
	Add(ctx context.Context, status *object.Status) (object.StatusID, error)
	// Delete Status
	DeleteByID(ctx context.Context, id object.StatusID) error
	// Find PublicTimeline
	FindPublicTimelines(ctx context.Context, onlyMedia bool, maxID int64, sinceID int64, limit int64) (object.Timelines, error)
}
