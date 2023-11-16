package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/ShayanGsh/azar/azar-db"
	"github.com/ShayanGsh/azar/internal/api"
	"github.com/ShayanGsh/azar/internal/models"
	"github.com/ShayanGsh/azar/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

var first_user = models.UserData{
	Username: "test_user1",
	Password: "test1234",
	Email: "testing@gmail.com",
}

var token string

func TestAddUser(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	err = models.AddUserWithHash(q, first_user, ctx)

	assert.NoError(t, err)
}


func TestGetUserByUsername(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := models.UserData{
		Username: "test_user1",
	}

    result, err := models.GetUser(q, user, ctx)
    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }

    if result.Username != user.Username {
        t.Errorf("expected username %s, but got %s", user.Username, result.Username)
    }
	email := pgtype.Text{String: "testing@gmail.com"}
    if result.Email.String != email.String {
        t.Errorf("expected email %s, but got %s", email.String, result.Email.String)
    }

	if utils.CheckPasswordHash("test1234", result.Password) != true {
		t.Errorf("passwords do not match")
	}
}

func TestGetUserByEmail(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := models.UserData{
		Email: "testing@gmail.com",
	}

    result, err := models.GetUser(q, user, ctx)
    if err != nil {
        t.Fatalf("error getting user: %v", err)
    }
	username := "test_user1"
    if result.Username != username {
        t.Errorf("expected username %s, but got %s", user.Username, result.Username)
    }
	email := pgtype.Text{String: user.Email}
    if result.Email.String != email.String {
        t.Errorf("expected email %s, but got %s", email.String, result.Email.String)
    }
	if utils.CheckPasswordHash("test1234", result.Password) != true {
		t.Errorf("passwords do not match")
	}
}

func TestGetUserNoUsernameOrEmail(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := models.UserData{}

	_, err = models.GetUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error getting user, but got nil")
	}
}

func TestGetUserNotFound(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := models.UserData{
		Username: "test_user2",
	}

	_, err = models.GetUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error getting user, but got nil")
	}
}

func TestUpdateUser(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := models.UpdateUserData{
		Username: "test_user1",
		NewEmail: "testing1@gmail.com",
		NewPassword: "test12345",
	}

	err = models.UpdateUser(q, user, ctx)
	if err != nil {
		t.Fatalf("error updating user: %v", err)
	}

	result, err := models.GetUser(q, models.UserData{Username: user.Username}, ctx)
	if err != nil {
		t.Fatalf("error getting user: %v", err)
	}

	if result.Email.String != user.NewEmail {
		t.Errorf("expected email %s, but got %s", user.NewEmail, result.Email.String)
	}

	if utils.CheckPasswordHash(user.NewPassword, result.Password) != true {
		t.Errorf("passwords do not match")
	}
}

func TestUpdateUserNoUsernameOrEmail(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := models.UpdateUserData{}

	err = models.UpdateUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error updating user, but got nil")
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := models.UpdateUserData{
		Username: "test_user2",
		Email: "testing1@gmail.com",
	}

	err = models.UpdateUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error updating user, but got nil")
	}
}

func TestDeleteUser(t *testing.T) {
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := models.UserData{
		Username: "test_user1",
	}

	err = models.DeleteUser(q, user, ctx)
	if err != nil {
		t.Fatalf("error deleting user: %v", err)
	}

	_, err = models.GetUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error getting user, but got nil")
	}
}

func TestRegisterHandler(t *testing.T) {
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

	expected := api.ReplyMessage{
		Success: true,
		Message: "success",
		Status:  http.StatusCreated,
	}

	parsed := api.ReplyMessage{}
	json.Unmarshal(rr.Body.Bytes(), &parsed)

	if parsed != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v",
			parsed, expected)
	}
}

func TestPasswordFailRegisterHandler(t *testing.T) {
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

	expected := api.ReplyMessage{
		Success: false,
		Message: "password must be at least 8 characters long",
		Status:  http.StatusInternalServerError,
	}

	parsed := api.ReplyMessage{}
	json.Unmarshal(rr.Body.Bytes(), &parsed)

	if parsed != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v",
			parsed, expected)
	}
}

func TestLoginHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBufferString(`{"username":"test","password":"Testing123"}`))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(uc.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v",
			status, http.StatusOK)
	}

	expected := api.ReplyMessage{
		Success: true,
		Message: "",
		Status:  http.StatusOK,
	}

	parsed := api.ReplyMessage{}
	json.Unmarshal(rr.Body.Bytes(), &parsed)

	token = parsed.Message

	if parsed.Status != expected.Status {
		t.Errorf("handler returned unexpected body: got %v, want %v",
			parsed, expected)
	}
}
