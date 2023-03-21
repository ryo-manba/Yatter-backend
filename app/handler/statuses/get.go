package statuses

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for `GET /v1/statuses/{id}`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	statusRepo := h.app.Dao.Status() // domain/repository の取得
	status, err := statusRepo.FindWithAccountByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	log.Println(fmt.Sprintf("Status: %+v", status))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
