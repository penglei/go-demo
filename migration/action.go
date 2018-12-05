package migration

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/qcloud2018/go-demo/config"
	"go.uber.org/zap"
	"net/url"
)

type Logger struct{}

func (_ *Logger) Printf(format string, v ...interface{}) {
	var l = zap.L().Sugar()
	l.Infof(format, v)
}

func (_ *Logger) Verbose() bool {
	return true
}

type FileMigration struct {
	Host      string
	Port      uint
	User      string
	Password  string
	Name      string
	Charset   string
	Versions  string
	SourceUrl string
}

func NewFileMigration(c config.Database) *FileMigration {
	return &FileMigration{
		Host:     c.Host,
		Port:     c.Port,
		User:     c.User,
		Password: c.Password,
		Name:     c.Name,
		Charset:  c.Charset,
	}
}

func (m *FileMigration) create() *migrate.Migrate {
	var l = zap.L()
	var sourceUrl = fmt.Sprintf("file://%s", m.Versions)
	var dbUrl string
	var userInfo = url.UserPassword(m.User, m.Password).String()
	dbUrl = fmt.Sprintf("mysql://%s@tcp(%s:%d)/%s?charset=%s",
		userInfo, m.Host, m.Port, m.Name, m.Charset)

	l.Info("migrate from versions in directory", zap.String("sourceUrl", sourceUrl))
	l.Info("migrate target database", zap.String("dbUrl", dbUrl))
	migrateInstance, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		l.Fatal("create migrate connection failed", zap.Error(err))
	}
	migrateInstance.Log = &Logger{}
	m.SourceUrl = sourceUrl
	return migrateInstance
}

func (m *FileMigration) Upgrade() {
	var l = zap.L()
	l.Info("do upgrade...")
	var migrateInstance = m.create()
	var err = migrateInstance.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			l.Warn("migrate.ErrNoChange", zap.Error(err))
		} else if _, ok := err.(migrate.ErrDirty); ok {
			l.Error("Last migration failed, you must solve the database state problem and try again！")
		} else {
			l.Error("Failed", zap.Error(err))
		}
	} else {
		l.Info("OK")
	}
}

func (m *FileMigration) Downgrade() {
	var l = zap.L()
	l.Info("do downgrade...")
	var migrateInstance = m.create()
	var err = migrateInstance.Steps(-1)
	if err != nil {
		l.Error("Failed", zap.Error(err))
	} else if _, ok := err.(migrate.ErrDirty); ok {
		l.Error("Last migration failed, you must solve the database state problem and try again！")
	} else {
		l.Info("OK")
	}
}

func (m *FileMigration) ForceResetDown() {
	var l = zap.L()
	var migrateInstance = m.create()
	var ver, dirty, err = migrateInstance.Version()
	if err != nil {
		l.Fatal("Get current version failed", zap.Error(err))
	}
	if !dirty {
		l.Fatal("Forbidden reset if database schema is not dirty.")
	}

	l.Info("Remove current version dirty state", zap.Uint("version", ver))
	err = migrateInstance.Force(int(ver))
	if err != nil {
		l.Fatal("Force version Failed", zap.Error(err))
	}

	l.Info("Downgrade current version ", zap.Uint("to_version", ver))
	err = migrateInstance.Steps(-1)
	if err != nil {
		l.Fatal("Force downgrade failed. You must solve the problem manually", zap.Error(err))
	} else {
		ver, _, _ = migrateInstance.Version()
		l.Info("Force downgrade successfully", zap.Uint("version", ver))
	}
}
