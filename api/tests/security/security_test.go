package security_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/Klaushayan/azar/api"
	"github.com/Klaushayan/azar/api/controllers"
	"github.com/Klaushayan/azar/api/tests/utils"
	db "github.com/Klaushayan/azar/azar-db"
)

// TODO: format this file

var ctx = context.Background()
var uc *controllers.UserController

func TestMain(m *testing.M) {
	// setup
	postgresContainer, mappedPort := test_utils.RunPostgresContainer()
	defer postgresContainer.Terminate(ctx)
	c := api.NewConfig()
	c.Database.Name = "azar_test"
	port, err := strconv.Atoi(mappedPort)
	if err != nil {
		log.Fatal(err)
	}
	c.Database.Port = port

	s := test_utils.GetServer(mappedPort, *c)
	s.MigrationCheck()
	uc = controllers.NewUserController(s.DB, s.JWTAuth)

	// run tests
	exitCode := m.Run()

	// teardown
	os.Exit(exitCode)

	m.Run()
}

func TestSQLInjection(t *testing.T) {

	// TODO: refactor the request/response code into a helper function
	rreq, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	body := `{"username":"test","password":"Testing123' OR 1=1; DROP TABLE users; --"}`

	rreq.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	rhandler := http.HandlerFunc(uc.Register)
	rhandler.ServeHTTP(rr, rreq)

	lreq, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	lreq.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	rr = httptest.NewRecorder()
	lhandler := http.HandlerFunc(uc.Login)
	lhandler.ServeHTTP(rr, lreq)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v",
			status, http.StatusOK)
	}

	// We should double check that the users table still exists

	ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := controllers.User{
		Username: "test",
		Password: "Testing123",
	}

	if _, err = uc.GetUser(q, user, ctx); err != nil {
		t.Fatal(err) // if the users table is dropped, this will fail
	}
}