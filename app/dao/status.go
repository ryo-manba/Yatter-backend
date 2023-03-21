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
	status struct {
		db *sqlx.DB
	}
)

// Create status repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// FindWIthAccountByID : アカウントの情報と共にステータスを取得する
func (r *status) FindWithAccountByID(ctx context.Context, id int64) (*object.Status, error) {
	query := `
	SELECT s.id,
				 s.content,
				 s.create_at as status_create_at,
				 a.id,
				 a.username,
				 a.password_hash,
				 a.display_name,
				 a.avatar,
				 a.header,
				 a.note,
				 a.create_at as account_create_at
	FROM status s
	INNER JOIN account a ON s.account_id = a.id
	WHERE s.id = ?
	`
	statusEntity := new(object.Status)
	accountEntity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, query, id).Scan(
		&statusEntity.ID,
		&statusEntity.Content,
		&statusEntity.CreateAt,
		&accountEntity.ID,
		&accountEntity.Username,
		&accountEntity.PasswordHash,
		&accountEntity.DisplayName,
		&accountEntity.Avatar,
		&accountEntity.Header,
		&accountEntity.Note,
		&accountEntity.CreateAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	statusEntity.Account = accountEntity

	return statusEntity, nil
}

// Add : 新規ステータス作成
func (r *status) Add(ctx context.Context, status *object.Status) (*object.Status, error) {
	query := `
	INSERT INTO status (account_id, content)
	VALUES (?, ?)
`
	println("status.AccountID", status.Account.ID)
	println("status.Content", status.Content)
	result, err := r.db.ExecContext(ctx, query, status.Account.ID, status.Content)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	status.ID = id
	// アカウントのデータとともにステータスを返却する
	return status, nil
}
