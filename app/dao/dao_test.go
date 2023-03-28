package dao

import (
	"context"
	"testing"
	"time"
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

	expected := &object.Account{
		Username:     "testuser",
		PasswordHash: "passwordhash",
		DisplayName:  toPtr("Test User"),
		Note:         toPtr("Hello, world!"),
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

	account := &object.Account{
		Username:     "testuser",
		PasswordHash: "passwordhash",
		DisplayName:  toPtr("Test User"),
		Note:         toPtr("Hello, world!"),
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

// Status
func TestStatus_FindWithAccountByID(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	ctx := context.Background()

	// 結果となる値を用意する
	expectedStatus := &object.Status{
		ID:      1,
		Content: "Hello, world!",
		Account: &object.Account{
			ID:           1,
			Username:     "testuser",
			PasswordHash: "passwordhash",
			DisplayName:  toPtr("Test User"),
			Note:         toPtr("Hello, world!"),
		},
	}

	statusCreatedAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	accountCreatedAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	// クエリ結果として返されるモック行をセットアップする
	rows := sqlmock.NewRows([]string{"s.id", "s.content", "status_create_at", "a.id", "a.username", "a.password_hash", "a.display_name", "a.avatar", "a.header", "a.note", "account_create_at"}).
		AddRow(expectedStatus.ID, expectedStatus.Content, statusCreatedAt, expectedStatus.Account.ID, expectedStatus.Account.Username, expectedStatus.Account.PasswordHash, *expectedStatus.Account.DisplayName, expectedStatus.Account.Avatar, expectedStatus.Account.Header, *expectedStatus.Account.Note, accountCreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM status s INNER JOIN account a ON s.account_id = a.id WHERE s.id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	statusRepo := NewStatus(db)

	status, err := statusRepo.FindWithAccountByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, expectedStatus.ID, status.ID)
	assert.Equal(t, expectedStatus.Content, status.Content)
	assert.Equal(t, expectedStatus.Account.ID, status.Account.ID)
	assert.Equal(t, expectedStatus.Account.Username, status.Account.Username)
}

func TestStatus_Add(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	ctx := context.Background()

	expectedStatus := &object.Status{
		ID:      1,
		Content: "Hello, world!",
		Account: &object.Account{
			ID: 1,
		},
	}

	mock.ExpectExec("INSERT INTO status \\(account_id, content\\) VALUES \\(\\?, \\?\\)").
		WithArgs(expectedStatus.Account.ID, expectedStatus.Content).
		WillReturnResult(sqlmock.NewResult(expectedStatus.ID, 1))

	statusRepo := NewStatus(db)

	status := &object.Status{
		Account: expectedStatus.Account,
		Content: expectedStatus.Content,
	}

	savedStatus, err := statusRepo.Add(ctx, status)
	assert.NoError(t, err)
	assert.NotNil(t, savedStatus)
	assert.Equal(t, expectedStatus.ID, savedStatus.ID)
	assert.Equal(t, expectedStatus.Content, savedStatus.Content)
	assert.Equal(t, expectedStatus.Account.ID, savedStatus.Account.ID)
}

func TestStatus_DeleteByID(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	ctx := context.Background()

	mock.ExpectExec("DELETE FROM status WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	statusRepo := NewStatus(db)

	err := statusRepo.DeleteByID(ctx, 1)
	assert.NoError(t, err)
}

func TestTimeline_FindPublic(t *testing.T) {
	db, mock := setup(t)
	defer db.Close()

	ctx := context.Background()
	timelineRepo := NewTimeline(db)

	expectedStatus := &object.Status{
		ID:      1,
		Content: "Hello, world!",
		Account: &object.Account{
			ID:           1,
			Username:     "testuser",
			PasswordHash: "passwordhash",
			DisplayName:  toPtr("Test User"),
			Note:         toPtr("Hello, world!"),
		},
	}
	statusCreatedAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	accountCreatedAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	rows := sqlmock.NewRows([]string{"s.id", "s.content", "status_create_at", "a.id", "a.username", "a.password_hash", "a.display_name", "a.avatar", "a.header", "a.note", "account_create_at"}).
		AddRow(expectedStatus.ID, expectedStatus.Content, statusCreatedAt, expectedStatus.Account.ID, expectedStatus.Account.Username, expectedStatus.Account.PasswordHash, *expectedStatus.Account.DisplayName, expectedStatus.Account.Avatar, expectedStatus.Account.Header, *expectedStatus.Account.Note, accountCreatedAt)
	mock.ExpectQuery("^SELECT (.+) FROM status s INNER JOIN account a ON s.account_id = a.id").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	statuses, err := timelineRepo.FindPublic(ctx, false, 0, 0, 40)
	assert.NoError(t, err)
	assert.NotNil(t, statuses)
	assert.Len(t, statuses, 1)

	status := statuses[0]
	assert.Equal(t, expectedStatus.ID, status.ID)
	assert.Equal(t, expectedStatus.Content, status.Content)
	assert.Equal(t, expectedStatus.Account.ID, status.Account.ID)
	assert.Equal(t, expectedStatus.Account.Username, status.Account.Username)
	assert.Equal(t, expectedStatus.Account.PasswordHash, status.Account.PasswordHash)
	assert.Equal(t, *expectedStatus.Account.DisplayName, *status.Account.DisplayName)
	assert.Equal(t, expectedStatus.Account.Avatar, status.Account.Avatar)
	assert.Equal(t, expectedStatus.Account.Header, status.Account.Header)
	assert.Equal(t, *expectedStatus.Account.Note, *status.Account.Note)
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

func toPtr(s string) *string {
	return &s
}
