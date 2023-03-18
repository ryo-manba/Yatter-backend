package accounts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/accounts`
type AddRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	// モデルを生成する
	account := new(object.Account)
	account.Username = req.Username
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	log.Println(fmt.Sprintf("Account: %+v", account))

	accountRepo := h.app.Dao.Account() // domain/repository の取得
	addedAccount, err := accountRepo.Add(ctx, account)
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
