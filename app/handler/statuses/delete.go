package statuses

import (
	"encoding/json"
	"errors"
	"net/http"

	"yatter-backend-go/app/domain/customerror"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `DELETE /v1/statuses/{id}`
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	statusRepo := h.app.Dao.Status() // domain/repository の取得
	if err := statusRepo.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, customerror.ErrNotFound) {
			httperror.NotFound(w, err)
		} else {
			httperror.InternalServerError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&struct{}{}); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
