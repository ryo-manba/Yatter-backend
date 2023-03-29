package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"yatter-backend-go/app/domain/customerror"
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
func (r *status) FindWithAccountByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
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
func (r *status) Add(ctx context.Context, status *object.Status) (object.StatusID, error) {
	query := `
	INSERT INTO status (account_id, content)
	VALUES (?, ?)
`
	result, err := r.db.ExecContext(ctx, query, status.Account.ID, status.Content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteByID : ステータスの削除
func (r *status) DeleteByID(ctx context.Context, id object.StatusID) error {
	query := `
		DELETE FROM status
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// 削除したかを確かめる
	if affectedRows == 0 {
		return customerror.ErrNotFound
	}
	return nil
}

// FindPublic : 公開中のタイムラインを取得する
func (r *status) FindPublicTimelines(ctx context.Context, onlyMedia bool, maxID int64, sinceID int64, limit int64) (object.Timelines, error) {
	query := `
	SELECT s.id as status_id,
				 s.content,
				 s.create_at as status_create_at,
				 a.id as account_id,
				 a.username,
				 a.password_hash,
				 a.display_name,
				 a.avatar,
				 a.header,
				 a.note,
				 a.create_at as account_create_at
	FROM status s
	INNER JOIN account a ON s.account_id = a.id
	`
	// クエリの条件を格納する変数を用意
	whereClauses := make([]string, 0)
	args := make([]interface{}, 0)

	if onlyMedia {
		// NOTE: bonusでmediaのテーブルを追加する
		// whereClauses = append(whereClauses, "s.id IN (SELECT status_id FROM media)")
	}

	if maxID > 0 {
		whereClauses = append(whereClauses, "s.id <= ?")
		args = append(args, maxID)
	}

	if sinceID > 0 {
		whereClauses = append(whereClauses, "s.id >= ?")
		args = append(args, sinceID)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY s.create_at DESC"

	if limit <= 0 || limit > 80 {
		limit = 40
	}
	query += " LIMIT ?"

	args = append(args, limit)
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timelines := object.Timelines{}
	for rows.Next() {
		status := new(object.Status)
		account := new(object.Account)
		err := rows.Scan(
			&status.ID,
			&status.Content,
			&status.CreateAt,
			&account.ID,
			&account.Username,
			&account.PasswordHash,
			&account.DisplayName,
			&account.Avatar,
			&account.Header,
			&account.Note,
			&account.CreateAt,
		)
		if err != nil {
			return nil, err
		}
		status.Account = account
		timelines = append(timelines, *status)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("%w", err)
		}
	}

	return timelines, nil
}
