package tests

import (
	"context"
	"testing"

	"github.com/Klaushayan/azar/api/controllers"
	db "github.com/Klaushayan/azar/azar-db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

var first_user = controllers.User{
	Username: "test_user1",
	Password: "test1234",
	Email: "testing@gmail.com",
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	err = uc.AddUserWithHash(q, first_user, ctx)

	assert.NoError(t, err)
}


func TestGetUserByUsername(t *testing.T) {
    ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := controllers.User{
		Username: "test_user1",
	}

    result, err := uc.GetUser(q, user, ctx)
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

	if controllers.CheckPasswordHash("test1234", result.Password) != true {
		t.Errorf("passwords do not match")
	}
}

func TestGetUserByEmail(t *testing.T) {
    ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := controllers.User{
		Email: "testing@gmail.com",
	}

    result, err := uc.GetUser(q, user, ctx)
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
	if controllers.CheckPasswordHash("test1234", result.Password) != true {
		t.Errorf("passwords do not match")
	}
}

func TestGetUserNoUsernameOrEmail(t *testing.T) {
    ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := controllers.User{}

	_, err = uc.GetUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error getting user, but got nil")
	}
}

func TestGetUserNotFound(t *testing.T) {
	ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)

	user := controllers.User{
		Username: "test_user2",
	}

	_, err = uc.GetUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error getting user, but got nil")
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := controllers.UpdateUser{
		Username: "test_user1",
		NewEmail: "testing1@gmail.com",
		NewPassword: "test12345",
	}

	err = uc.UpdateUser(q, user, ctx)
	if err != nil {
		t.Fatalf("error updating user: %v", err)
	}

	result, err := uc.GetUser(q, controllers.User{Username: user.Username}, ctx)
	if err != nil {
		t.Fatalf("error getting user: %v", err)
	}

	if result.Email.String != user.NewEmail {
		t.Errorf("expected email %s, but got %s", user.NewEmail, result.Email.String)
	}

	if controllers.CheckPasswordHash(user.NewPassword, result.Password) != true {
		t.Errorf("passwords do not match")
	}
}

func TestUpdateUserNoUsernameOrEmail(t *testing.T) {
	ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := controllers.UpdateUser{}

	err = uc.UpdateUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error updating user, but got nil")
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := controllers.UpdateUser{
		Username: "test_user2",
		Email: "testing1@gmail.com",
	}

	err = uc.UpdateUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error updating user, but got nil")
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	c, err := uc.DatabaseConnectionPool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	q := db.New(c)
	user := controllers.User{
		Username: "test_user1",
	}

	err = uc.DeleteUser(q, user, ctx)
	if err != nil {
		t.Fatalf("error deleting user: %v", err)
	}

	_, err = uc.GetUser(q, user, ctx)
	if err == nil {
		t.Fatalf("expected error getting user, but got nil")
	}
}