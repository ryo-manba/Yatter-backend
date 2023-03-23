package app_test

import (
	"testing"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/config"
	"yatter-backend-go/app/dao"
)

func setupTestDB(t *testing.T) (*app.App, func()) {
	testCfg := config.MySQLTestConfig()
	testDao, err := dao.New(testCfg)
	if err != nil {
		t.Fatalf("Failed to create testDao: %v", err)
	}
	// TODO: テスト用のデータを入れる

	appInstance, err := app.NewAppWithDao(testDao)
	if err != nil {
		t.Fatalf("Failed to create appInstance: %v", err)
	}

	cleanup := func() {
		// TODO: データを削除する
	}

	return appInstance, cleanup
}

func TestInit(t *testing.T) {
	println("Test started")
	_, cleanup := setupTestDB(t)
	defer cleanup()
	println("Test ended")
}
