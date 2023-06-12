package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Klaushayan/azar/api"
	"github.com/Klaushayan/azar/api/controllers"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ctx = context.Background()
var uc *controllers.UserController

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
		time.Sleep(2 * time.Second)
		s.Shutdown()
	}()

	time.Sleep(3 * time.Second)

	if s.IsRunning() {
		t.Fatal("server should be stopped")
	}
}

func TestRegisterHandler(t *testing.T) {
	s := getServer()
	uc = controllers.NewUserController(s.DB, s.JWTAuth)
    req, err := http.NewRequest("GET", "/register", nil)
    if err != nil {
        t.Fatal(err)
    }
	req.Body = ioutil.NopCloser(bytes.NewBufferString(`{"username":"test","password":"Testing123"}`))
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(uc.Register)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v, want %v",
            status, http.StatusCreated)
    }

	expected := controllers.ReplyMessage{
		Success: true,
		Message: "success",
		Status:  http.StatusCreated,
		Error:   nil,
	}

	parsed := controllers.ReplyMessage{}
	json.Unmarshal(rr.Body.Bytes(), &parsed)

	if parsed != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v",
			parsed, expected)
	}
}

func TestPasswordFailRegisterHandler(t *testing.T) {
	s := getServer()
	uc = controllers.NewUserController(s.DB, s.JWTAuth)
    req, err := http.NewRequest("GET", "/register", nil)
    if err != nil {
        t.Fatal(err)
    }
	req.Body = ioutil.NopCloser(bytes.NewBufferString(`{"username":"test1","password":"test"}`))
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(uc.Register)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("handler returned wrong status code: got %v, want %v",
            status, http.StatusInternalServerError)
    }

	expected := controllers.ReplyMessage{
		Success: false,
		Message: "password must be at least 8 characters long",
		Status:  http.StatusInternalServerError,
		Error:   nil,
	}

	parsed := controllers.ReplyMessage{}
	json.Unmarshal(rr.Body.Bytes(), &parsed)

	if parsed != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v",
			parsed, expected)
	}
}