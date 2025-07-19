package user

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func setup(ctx context.Context, t *testing.T) *pgx.Conn {
	conn, err := pgx.Connect(ctx, "postgres://admin:admin@127.0.0.1:5432/chatdb?sslmode=disable")
	if err != nil {
		t.Fatal("DB connection err:", err)
	}

	sqlBytes, err := os.ReadFile("../database/001_add_user_table.up.sql")
	if err != nil {
		t.Error("Read sql failed, err:", err)
	}
	if _, err := conn.Exec(ctx, string(sqlBytes)); err != nil {
		t.Error("Create user table failed, err:", err)
	}

	return conn
}

func clean(ctx context.Context, t *testing.T, conn *pgx.Conn) {
	defer conn.Close(ctx)

	sqlBytes, err := os.ReadFile("../database/001_add_user_table.down.sql")
	if err != nil {
		t.Error("Read sql failed, err:", err)
	}
	if _, err := conn.Exec(ctx, string(sqlBytes)); err != nil {
		t.Error("Delete user table failed, err:", err)
	}
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	conn := setup(ctx, t)
	defer clean(ctx, t, conn)

	repo := &UserRepo{conn: conn}
	pa := CreateParams{
		Username: "USER1",
		Email:    "user1@email.com",
		Password: "password",
	}
	user, err := repo.CreateUser(ctx, &pa)
	if err != nil {
		t.Error("err:", err)
		t.FailNow()
	}
	if user.Email != pa.Email || user.Username != pa.Username {
		t.Error("email or username doesn't match")
	}
	if string(user.PasswordHash) == string(pa.Password) {
		t.Error("plain password is stored directly in DB")
	}
	if !CheckPassword(pa.Password, user.PasswordHash) {
		t.Error("hash password logic failed")
	}
}
