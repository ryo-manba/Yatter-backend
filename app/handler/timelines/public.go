package timelines

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Request body for `GET /v1/timelines/public`
type Params struct {
	OnlyMedia bool
	MaxID     object.AccountID
	SinceID   object.AccountID
	Limit     int64
}

const (
	DefaultOnlyMedia = false
	DefaultMaxID     = 0
	DefaultSinceID   = 0
	DefaultLimit     = 40
	MaxLimit         = 80
)

// Handle request for `GET /v1/timelines/public`
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := parse(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	timelineRepo := h.app.Dao.Timeline() // domain/repository の取得

	timeline, err := timelineRepo.FindPublic(ctx, params.OnlyMedia, params.MaxID, params.SinceID, params.Limit)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// Userの情報を返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timeline); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func parse(r *http.Request) (*Params, error) {
	onlyMedia, err := request.QueryBool(r, "only_media", false)
	if err != nil {
		return nil, err
	}
	maxID, err := request.QueryInt64(r, "max_id", DefaultMaxID)
	if err != nil {
		return nil, err
	}
	sinceID, err := request.QueryInt64(r, "since_id", DefaultSinceID)
	if err != nil {
		return nil, err
	}
	limit, err := request.QueryInt64(r, "limit", DefaultLimit)
	if err != nil {
		return nil, err
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return &Params{
		OnlyMedia: onlyMedia,
		MaxID:     maxID,
		SinceID:   sinceID,
		Limit:     limit,
	}, nil
}
