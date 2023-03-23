package timelines

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Status   string `json:"status"`
	MediaIds []int  `json:"media_ids"`
}

// Handle request for `GET /v1/timelines/public`
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: エラーチェックをする
	onlyMedia := r.URL.Query().Get("only_media") == "1"
	maxID, _ := strconv.ParseInt(r.URL.Query().Get("max_id"), 10, 64)
	sinceID, _ := strconv.ParseInt(r.URL.Query().Get("since_id"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)

	timelineRepo := h.app.Dao.Timeline() // domain/repository の取得

	timeline, err := timelineRepo.FindPublicTimeline(ctx, onlyMedia, maxID, sinceID, limit)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	log.Println(fmt.Sprintf("Account: %+v", timeline))
	// Userの情報を返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(timeline); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
