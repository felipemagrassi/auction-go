package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/felipemagrassi/auction-go/configuration/logger"
	"github.com/felipemagrassi/auction-go/internal/entity/user_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRepository) FindById(ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": userId}

	var user UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			message := fmt.Sprintf("User not found with this id = %s", userId)
			logger.Error(message, err)
			return nil, internal_error.NewNotFoundError(
				message,
			)
		}
		message := fmt.Sprintf("Error trying to find user with id = %s", userId)
		logger.Error(message, err)
		return nil, internal_error.NewInternalServerError(
			message,
		)
	}

	userEntity := &user_entity.User{
		Id:   user.Id,
		Name: user.Name,
	}

	return userEntity, nil
}
