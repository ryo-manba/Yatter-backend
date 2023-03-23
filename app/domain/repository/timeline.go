package repository

import (
	"context"
	"yatter-backend-go/app/domain/object"
)

type Timeline interface {
	FindPublicTimeline(ctx context.Context, onlyMedia bool, maxID int64, sinceID int64, limit int64) ([]*object.Status, error)
}
