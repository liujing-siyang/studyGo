package biz

import (
	"context"
	"errors"
	"realworld/internal/conf"
	"realworld/internal/pkg/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email        string
	Token        string
	Username     string
	Bio          string
	Image        string
	PasswordHash string
}

type UserLogin struct {
	Email    string
	Token    string
	Username string
	Bio      string
	Image    string
}

func hashPassword(pwd string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func verifyPassword(hashPassword, inputPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(inputPassword)); err != nil {
		return false
	}
	return true
}

type UserRepo interface {
	CreatUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type ProfileRepo interface {
}

type UserUsecase struct {
	ur  UserRepo
	pr  ProfileRepo
	jwt *conf.Jwt
	log *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewUserUsecase(ur UserRepo,
	pr ProfileRepo,
	jwt *conf.Jwt, logger log.Logger) *UserUsecase {
	return &UserUsecase{ur: ur, pr: pr, jwt: jwt, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) generateToken(username string) string {
	return auth.GenerateToken(uc.jwt.Secret, username)
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *UserUsecase) Register(ctx context.Context, email, username, password string) (*UserLogin, error) {
	u := &User{
		Email:        email,
		Username:     username,
		PasswordHash: hashPassword(password),
	}
	if err := uc.ur.CreatUser(ctx, u); err != nil {
		return nil, err
	}
	return &UserLogin{
		Email:    email,
		Username: username,
		Token:    uc.generateToken(u.Username),
	}, nil
}

func (uc *UserUsecase) Login(ctx context.Context, email string, password string) (*UserLogin, error) {
	u, err := uc.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !verifyPassword(u.PasswordHash, password) {
		return nil, errors.New("Login failed")
	}
	return &UserLogin{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    u.Token,
	}, nil
}
