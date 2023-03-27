package dao

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// Account
func TestAccount_FindByUsername(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	ctx := context.Background()

	// 結果となる値をセットアップする
	displayName := "Test User"
	note := "Hello, world!"
	expected := &object.Account{
		Username:     "testuser",
		PasswordHash: "passwordhash",
		DisplayName:  &displayName,
		Note:         &note,
	}

	// クエリ結果として返されるモック行をセットアップする
	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "display_name", "avatar", "header", "note"}).
		AddRow(1, expected.Username, expected.PasswordHash, expected.DisplayName, expected.Avatar, expected.Header, expected.Note)

	// クエリとその引数の期待値を設定する
	mock.ExpectQuery("(?i)SELECT (.+) FROM account WHERE username = ?").
		WithArgs(expected.Username).
		WillReturnRows(rows)

	accountRepo := NewAccount(db)
	account, err := accountRepo.FindByUsername(ctx, expected.Username)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, int64(1), account.ID)
	assert.Equal(t, expected.Username, account.Username)
	assert.Equal(t, expected.PasswordHash, account.PasswordHash)
	assert.Equal(t, *expected.DisplayName, *account.DisplayName)
	assert.Equal(t, expected.Avatar, account.Avatar)
	assert.Equal(t, expected.Header, account.Header)
	assert.Equal(t, *expected.Note, *account.Note)
}

func TestAccount_Add(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	ctx := context.Background()
	accountRepo := NewAccount(db)

	displayName := "Test User"
	note := "Hello, world!"

	account := &object.Account{
		Username:     "testuser",
		PasswordHash: "passwordhash",
		DisplayName:  &displayName,
		Note:         &note,
	}
	// Setup mock
	mock.ExpectExec("(?i)INSERT INTO account (.+) VALUES (.+)").
		WithArgs(account.Username, account.PasswordHash, account.DisplayName, account.Avatar, account.Header, account.Note).
		WillReturnResult(sqlmock.NewResult(1, 1))

	savedAccount, err := accountRepo.Add(ctx, account)
	assert.NoError(t, err)
	assert.NotNil(t, savedAccount)
	assert.Equal(t, int64(1), savedAccount.ID)
	assert.Equal(t, account.Username, savedAccount.Username)
	assert.Equal(t, account.PasswordHash, savedAccount.PasswordHash)
	assert.Equal(t, *account.DisplayName, *savedAccount.DisplayName)
	assert.Equal(t, account.Avatar, savedAccount.Avatar)
	assert.Equal(t, account.Header, savedAccount.Header)
	assert.Equal(t, *account.Note, *savedAccount.Note)
}

// Utils
func setup(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	rawDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	// sqlmock データベース接続を sqlx でラップする
	db := sqlx.NewDb(rawDb, "sqlmock")
	return db, mock
}
