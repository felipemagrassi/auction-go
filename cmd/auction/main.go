package main

import (
	"context"
	"log"

	"github.com/felipemagrassi/auction-go/configuration/database/mongodb"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
