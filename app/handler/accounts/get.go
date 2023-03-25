package accounts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username, err := request.UsernameOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	accountRepo := h.app.Dao.Account() // domain/repository の取得

	account, err := accountRepo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	log.Println(fmt.Sprintf("Account: %+v", account))
	// Userの情報を返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
