package timelines

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Status   string `json:"status"`
	MediaIds []int  `json:"media_ids"`
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

	onlyMedia, err := request.QueryBool(r, "only_media", false)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	maxID, err := request.QueryInt64(r, "max_id", DefaultMaxID)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	sinceID, err := request.QueryInt64(r, "since_id", DefaultSinceID)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	limit, err := request.QueryInt64(r, "limit", DefaultLimit)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	timelineRepo := h.app.Dao.Timeline() // domain/repository の取得

	timeline, err := timelineRepo.FindPublic(ctx, onlyMedia, maxID, sinceID, limit)
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
