package app

import (
	"yatter-backend-go/app/config"
	"yatter-backend-go/app/dao"
)

// Dependency manager for whole application
type App struct {
	Dao dao.Dao
}

// Create dependency manager
func NewApp() (*App, error) {
	// panic if lacking something
	daoCfg := config.MySQLConfig()

	dao, err := dao.New(daoCfg)
	if err != nil {
		return nil, err
	}

	return &App{Dao: dao}, nil
}

// Create dependency manager for tests
func NewTestApp() (*App, error) {
	testCfg := config.MySQLTestConfig()
	dao, err := dao.New(testCfg)
	if err != nil {
		return nil, err
	}

	return &App{Dao: dao}, nil
}
