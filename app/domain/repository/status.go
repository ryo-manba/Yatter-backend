package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Find Status
	FindWithAccountByID(ctx context.Context, id int64) (*object.Status, error)
	// Create Status
	Add(ctx context.Context, status *object.Status) (*object.Status, error)

	// Find Status
	//	FindById(ctx context.Context, id int64) (*object.Status, error)
}
