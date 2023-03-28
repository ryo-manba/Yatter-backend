package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	timeline struct {
		db *sqlx.DB
	}
)

// Create timeline repository
func NewTimeline(db *sqlx.DB) repository.Timeline {
	return &timeline{db: db}
}

// FindPublic : 公開中のタイムラインを取得する
// パラメータによって取得内容を変更する
func (r *timeline) FindPublic(ctx context.Context, onlyMedia bool, maxID int64, sinceID int64, limit int64) ([]*object.Status, error) {
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

	statuses := make([]*object.Status, 0)
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
		statuses = append(statuses, status)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("%w", err)
		}
	}

	return statuses, nil
}
