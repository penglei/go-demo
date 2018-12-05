package test

import (
	"database/sql"
	"github.com/qcloud2018/go-demo/service"
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

// SetupEnv creates a new test environment, including a clean database and an instance of our HTTP service.
func SetupEnv(t *testing.T) *Env {
	db := SetupDB(t)
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
func SetupDB(t *testing.T) *service.Database {
	var databaseURL string
	db, err := sql.Open("postgres", databaseURL)
	require.NoError(t, err, "Error opening database")

	return &service.Database{db}
}
