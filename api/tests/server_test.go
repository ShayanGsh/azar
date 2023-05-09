package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Klaushayan/azar/api"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ctx = context.Background()

// postgres

func TestMain(m *testing.M) {
	// setup
	postgresContainer := runPostgresContainer()
	defer postgresContainer.Terminate(ctx)

	// run tests
	exitCode := m.Run()

	// teardown
	os.Exit(exitCode)
}

var mappedPort nat.Port

func runPostgresContainer() testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13.2",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "azar_test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err = postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(err)
	}

	postgresHost, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	postgresDSN := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable",
	postgresHost, mappedPort.Port())

	os.Setenv("POSTGRES_DSN", postgresDSN)

	// wait for postgres to start
	time.Sleep(5 * time.Second)

	return postgresContainer
}

func getServer() *api.Server {
	c, err := api.LoadConfig("config_example.json")
	c.Database.Port = mappedPort.Int()
	if err != nil {
		log.Fatal(err)
	}
	s := api.NewServer(c)
	return s
}

func TestMigration(t *testing.T) {
	s := getServer()
	b := s.MigrationCheck()
	if !b {
		t.Fatal("migration failed")
	}
}

func TestStartingServer(t *testing.T) {
	s := getServer()

	go func() {
		s.Start()
	}()

	go func() {
		time.Sleep(5 * time.Second)
		s.Shutdown()
	}()

	time.Sleep(10 * time.Second)

	if s.IsRunning() {
		t.Fatal("server should be stopped")
	}
}