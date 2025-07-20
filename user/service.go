package user

import (
	"context"

	"local/chat"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", chat.NewError(chat.ERR_INTERNAL, err.Error())
	}
	return string(hashed), nil
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserService struct {
	repo *UserRepo
}

func NewUserService(conn *pgx.Conn) *UserService {
	return &UserService{repo: NewUserRepo(conn)}
}

func (s *UserService) Register(ctx context.Context, pa *RegisterParams) (User, error) {
	return s.repo.CreateUser(ctx, pa)
}
