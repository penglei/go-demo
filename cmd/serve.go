package main

import (
	"database/sql"
	"fmt"
	"github.com/qcloud2018/go-demo/config"
	"github.com/qcloud2018/go-demo/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run web server",
	Long:  `run http server which provide contact manage service`,
	Run: func(cmd *cobra.Command, args []string) {
		zap.L().Info("do serve command.")
		runServer()
	},
}

func runServer() {

	db := SetupDB(conf.Database)
	server := service.NewServer(db)
	http.HandleFunc("/", server.ServeHTTP)

	var err = http.ListenAndServe(":8080", nil)
	if err != nil {
		zap.L().Error("server exit error", zap.Error(err))
	}
}

func SetupDB(dbConf config.Database) *service.Database {
	databaseURL := dbConf.GetURL()
	if databaseURL == "" {
		panic("db config must be set!")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		panic(fmt.Sprintf("Unable to open DB connection: %+v", err))
	}

	return &service.Database{DB: db}
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
