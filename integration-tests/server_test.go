package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ShayanGsh/azar/api"
	"github.com/ShayanGsh/azar/api/controllers"
	"github.com/ShayanGsh/azar/integration-tests/utils"
)

var ctx = context.Background()
var uc *controllers.UserController
var s *api.APIServer

// postgres

func TestMain(m *testing.M) {
	// setup
	postgresContainer, mappedPort := test_utils.RunPostgresContainer()
	defer postgresContainer.Terminate(ctx)
	c := test_utils.GenerateConfig(mappedPort)
	s = api.NewAPIServer(c)
	uc = controllers.NewUserController(s.DB, s.JWTAuth)

	// run tests
	exitCode := m.Run()

	// teardown
	os.Exit(exitCode)
}

func TestMigration(t *testing.T) {
	b := s.MigrationCheck()
	if !b {
		t.Fatal("migration failed")
	}
}

func TestStartingServer(t *testing.T) {
	c := test_utils.GenerateConfig("5432")
	server := api.NewAPIServer(c)
	go func() {
		server.Start()
	}()

	go func() {
		time.Sleep(2 * time.Second)
		server.Shutdown()
	}()

	time.Sleep(3 * time.Second)

	if server.IsRunning() {
		t.Fatal("server should be stopped")
	}
}