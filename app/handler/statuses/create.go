package statuses

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Status   string `json:"status"`
	MediaIds []int  `json:"media_ids"`
}

// Handle request for `POST /v1/statuses`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	if err := req.Validate(); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status := new(object.Status)
	statusRepo := h.app.Dao.Status() // domain/repository の取得

	status.Content = req.Status

	// account の取得
	status.Account = auth.AccountOf(r)

	id, err := statusRepo.Add(ctx, status)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	addedStatus, _ := statusRepo.FindWithAccountByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// Userの情報を返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(addedStatus); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func parse(r *http.Request) (*AddRequest, error) {
	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Status == "" {
		return nil, fmt.Errorf("status not found")
	}
	return &req, nil
}

func (req *AddRequest) Validate() error {
	if req.Status == "" {
		return errors.New("status is required")
	}
	// bonus
	// if len(req.MediaIds) == 0 {
	// 	return errors.New("mediaID is required")
	// }
	return nil
}
