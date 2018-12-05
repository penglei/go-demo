package config

import "fmt"

type Database struct {
	Host     string
	Port     uint
	User     string
	Password string
	Name     string
}

func (db Database) GetURL() string {
	// "postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full"
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		db.User, db.Password, db.Host, db.Port, db.Name)
}

type Config struct {
	LogLevel  string // debug, info, warn, error, fatal or panic
	LogFormat string // json or console
	Env       string
	Addr      string
	Database  Database
}
