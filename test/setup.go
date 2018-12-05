package test

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql" //comment for lint
	"github.com/qcloud2018/go-demo/config"
	"github.com/qcloud2018/go-demo/logger"
	"github.com/qcloud2018/go-demo/service"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

// Env provides ess to all services used in tests, like the database, our server, and an HTTP client for performing
// HTTP requests against the test server.
type Env struct {
	T          *testing.T
	DB         *service.Database
	Server     *service.Server
	HTTPServer *httptest.Server
	Client     service.Client
}

// Close must be called after each test to ensure the Env is properly destroyed.
func (env *Env) Close() {
	env.HTTPServer.Close()
	env.DB.Close()
}

var conf = config.Config{}

func init() {
	flag.Parse()
	var cfgFile = flag.Arg(0)
	if cfgFile == "" {
		panic("no main config specified")
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO_DEMO")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	if err := logger.InitZapLogger(conf.Env, conf.LogLevel, conf.LogFormat); err != nil {
		panic(err)
	}

}

// SetupEnv creates a new test environment, including a clean database and an instance of our HTTP service.
func SetupEnv(t *testing.T) *Env {
	db := SetupDB(t, conf)
	server := service.NewServer(db)
	httpServer := httptest.NewServer(server)
	return &Env{
		T:          t,
		DB:         db,
		Server:     server,
		HTTPServer: httpServer,
		Client:     service.NewClient(httpServer.URL),
	}
}

// SetupDB initializes a test database, performing all migrations.
func SetupDB(t *testing.T, cfg config.Config) *service.Database {
	var databaseURL = cfg.Database.GetURL()
	db, err := sql.Open("mysql", databaseURL)
	require.NoError(t, err, "Error opening database")

	return &service.Database{DB: db}
}
