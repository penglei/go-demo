package test

import (
	"database/sql"
	"github.com/mattes/migrate/migrate"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"os"
	"testing"
	"workshop-demo/service"
)

// Env provides access to all services used in tests, like the database, our server, and an HTTP client for performing
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
	databaseURL := os.Getenv("CONTACTS_DB_URL")
	require.NotEmpty(t, databaseURL, "CONTACTS_DB_URL must be set!")

	sqlFiles := "./db/migrations"
	if sqlFilesEnv := os.Getenv("CONTACTS_DB_MIGRATIONS"); sqlFilesEnv != "" {
		sqlFiles = sqlFilesEnv
	}
	allErrors, ok := migrate.ResetSync(databaseURL, sqlFiles)
	require.True(t, ok, "Failed to migrate database %v", allErrors)

	db, err := sql.Open("postgres", databaseURL)
	require.NoError(t, err, "Error opening database")

	return &service.Database{db}
}
