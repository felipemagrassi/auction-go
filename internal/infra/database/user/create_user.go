package user

import (
	"context"
	"fmt"

	"github.com/felipemagrassi/auction-go/configuration/logger"
	"github.com/felipemagrassi/auction-go/internal/entity/user_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
)

func (ur *UserRepository) CreateUser(ctx context.Context, user *user_entity.User) *internal_error.InternalError {
	userEntity := UserEntityMongo{
		Id:   user.Id,
		Name: user.Name,
	}

	_, err := ur.Collection.InsertOne(ctx, userEntity)
	if err != nil {
		fmt.Printf(err.Error())
		message := fmt.Sprintf("Error trying to create user with id = %s", user.Id)
		logger.Error(message, err)
		return internal_error.NewInternalServerError(
			message,
		)
	}

	return nil
}
