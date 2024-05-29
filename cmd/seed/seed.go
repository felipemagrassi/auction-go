package main

import (
	"context"
	"log"

	"github.com/felipemagrassi/auction-go/configuration/database/mongodb"
	"github.com/felipemagrassi/auction-go/internal/entity/user_entity"
	"github.com/felipemagrassi/auction-go/internal/infra/database/user"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal(err.Error())
		return
	}

	database, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	userRepository := user.NewUserRepository(database)
	userEntity := &user_entity.User{Id: uuid.New().String(), Name: "Felipe"}

	if err := userRepository.CreateUser(ctx, userEntity); err != nil {
		log.Fatal(err.Error())
	}
}
