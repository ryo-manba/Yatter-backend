package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}

// Add : 新規ユーザ作成
func (r *account) Add(ctx context.Context, account *object.Account) (*object.Account, error) {
	query := `
		INSERT INTO account (username, password_hash, display_name, avatar, header, note)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		account.Username,
		account.PasswordHash,
		account.DisplayName,
		account.Avatar,
		account.Header,
		account.Note,
	)

	if err != nil {
		return nil, err
	}

	// 挿入に成功した場合にidを取得して、accountのIDに設定する
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	account.ID = id
	return account, nil
}
