package security_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ShayanGsh/azar/api"
	"github.com/ShayanGsh/azar/api/controllers"
	"github.com/ShayanGsh/azar/api/tests/utils"
	db "github.com/ShayanGsh/azar/azar-db"
	"github.com/ShayanGsh/azar/core"
)

// TODO: format this file

var ctx = context.Background()
var uc *controllers.UserController

func TestMain(m *testing.M) {
	// setup
	postgresContainer, mappedPort := test_utils.RunPostgresContainer()
	defer postgresContainer.Terminate(ctx)
	c := test_utils.GenerateConfig(mappedPort)
	s := api.NewAPIServer(c)

	s.MigrationCheck()
	uc = controllers.NewUserController(s.DB, s.JWTAuth)

	// run tests
	exitCode := m.Run()

	// teardown
	os.Exit(exitCode)
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

	user := core.UserData{
		Username: "test",
		Password: "Testing123",
	}

	if _, err = core.GetUser(q, user, ctx); err != nil {
		t.Fatal(err) // if the users table is dropped, this will fail
	}
}