package user

import (
	"context"
	"time"

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

func (r *UserRepo) CreateUser(ctx context.Context, p *RegisterParams) (User, error) {
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
