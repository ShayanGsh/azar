package test_utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Klaushayan/azar/api"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ctx = context.Background()


func RunPostgresContainer() (testcontainers.Container, string) {
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

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
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

	return postgresContainer, mappedPort.Port()
}

func GetServer(port string, config api.Config) *api.Server {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	config.Port = intPort
	s := api.NewServer(&config)
	return s
}