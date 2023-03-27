package dao

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// DAO interface
	Dao interface {
		// Get account repository
		Account() repository.Account

		// Get status repository
		Status() repository.Status

		// Get timeline repository
		Timeline() repository.Timeline

		// Clear all data in DB
		InitAll() error

		// Setup Test DB data
		SetupTestDB() error
	}

	// Implementation for DAO
	dao struct {
		db *sqlx.DB
	}
)

// Create DAO
func New(config DBConfig) (Dao, error) {
	db, err := initDb(config)
	if err != nil {
		return nil, err
	}

	return &dao{db: db}, nil
}

func (d *dao) Account() repository.Account {
	return NewAccount(d.db)
}

func (d *dao) Status() repository.Status {
	return NewStatus(d.db)
}

func (d *dao) Timeline() repository.Timeline {
	return NewTimeline(d.db)
}

// 外部キー制約を無効にしてから、テーブルを削除してる
func (d *dao) InitAll() error {
	if err := d.exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return fmt.Errorf("Can't disable FOREIGN_KEY_CHECKS: %w", err)
	}

	defer func() {
		err := d.exec("SET FOREIGN_KEY_CHECKS=0")
		if err != nil {
			log.Printf("Can't restore FOREIGN_KEY_CHECKS: %+v", err)
		}
	}()

	for _, table := range []string{"account", "status"} {
		if err := d.exec("TRUNCATE TABLE " + table); err != nil {
			return fmt.Errorf("Can't truncate table "+table+": %w", err)
		}
	}

	return nil
}

const seedPath = "/work/yatter-backend-go/ddl/tool/seed.sql"

// seedの値を入れる
func (d *dao) SetupTestDB() error {
	content, err := ioutil.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("Failed to read seed file : %+v", err)
	}
	// SQLステートメントを実行します。
	sqlStatements := strings.Split(string(content), ";")
	for _, stmt := range sqlStatements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		err := d.exec(stmt)
		if err != nil {
			log.Fatalf("Failed to execute statement: %v\n%v", stmt, err)
			return fmt.Errorf("Failed to execute seed query : %v\n%+v", stmt, err)
		} else {
			// FIXME: debug用
			// fmt.Printf("Successfully executed statement: %v\n", stmt)
		}
	}
	return nil
}

func (d *dao) exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)
	return err
}
