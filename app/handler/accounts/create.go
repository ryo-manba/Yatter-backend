package accounts

import (
	"encoding/json"
	"errors"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/accounts`
type AddRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Note        string `json:"note"`
	Avatar      string `json:"avatar"`
	Header      string `json:"header"`
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// パラメータのバリデーション
	if err := req.Validate(); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// モデルを生成する
	account := &object.Account{
		Username:    req.Username,
		Avatar:      &req.Avatar,
		Note:        &req.Note,
		Header:      &req.Header,
		DisplayName: &req.DisplayName,
	}
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	accountRepo := h.app.Dao.Account() // domain/repository の取得
	_, err := accountRepo.Add(ctx, account)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	addedAccount, err := accountRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// Userの情報を返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(addedAccount); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (req *AddRequest) Validate() error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
