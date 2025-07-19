package user

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type CreateParams struct {
	Email    string
	Username string
	Password string
}

type UserRepo struct {
	conn *pgx.Conn
}

func (r *UserRepo) CreateUser(ctx context.Context, p *CreateParams) (User, error) {
	sql := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, password_hash, created_at
	`
	var user User
	pwHash, err := HashPassword(p.Password)
	if err != nil {
		return user, err
	}
	row := r.conn.QueryRow(ctx, sql, p.Username, p.Email, pwHash)
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
