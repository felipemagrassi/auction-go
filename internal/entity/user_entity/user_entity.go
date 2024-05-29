package user_entity

import (
	"context"

	"github.com/felipemagrassi/auction-go/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(context.Context, string) (*User, *internal_error.InternalError)
}
