package user

import (
	"context"
	"errors"
	"local/chat"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type RegisterParams struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserRepo struct {
	conn *pgx.Conn
}

func NewUserRepo(conn *pgx.Conn) *UserRepo {
	return &UserRepo{conn: conn}
}

func (r *UserRepo) CreateUser(ctx context.Context, pa *RegisterParams) (User, error) {
	sql := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, password_hash, created_at
	`
	var user User

	pwHash, err := HashPassword(pa.Password)
	if err != nil {
		return user, err
	}
	err = r.conn.QueryRow(ctx, sql, pa.Username, pa.Email, pwHash).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)

	var pgError *pgconn.PgError
	if ok := errors.As(err, &pgError); ok {
		if pgError.Code == "23505" {
			return user, chat.NewError(chat.ERR_DUPLICATE, "duplicate username or email")
		}
	} else if err != nil {
		return user, err
	}

	return user, nil
}
