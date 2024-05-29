package auction

import (
	"context"
	"fmt"
	"time"

	"github.com/felipemagrassi/auction-go/configuration/logger"
	"github.com/felipemagrassi/auction-go/internal/entity/auction_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ar *AuctionRepository) FindAuctionById(
	ctx context.Context,
	auctionId string,
) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": auctionId}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		message := fmt.Sprintf("Error trying to find auction by id = %s", auctionId)
		logger.Error(message, err)
		return nil, internal_error.NewInternalServerError(message)
	}

	return &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Category:    auctionEntityMongo.Category,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
	}, nil
}

func (ar *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category string,
	productName string,
) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}
	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["product_name"] = primitive.Regex{
			Pattern: productName,
			Options: "i",
		}
	}

	var auctionEntitiesMongo []AuctionEntityMongo

	cursor, err := ar.Collection.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		message := fmt.Sprintf("Error trying to find auction by filter = %v", filter)
		logger.Error(message, err)
		return nil, internal_error.NewInternalServerError(message)
	}

	if err := cursor.All(ctx, &auctionEntitiesMongo); err != nil {
		message := fmt.Sprintf("Error parsing auctions by filter = %v", filter)
		logger.Error(message, err)
		return nil, internal_error.NewInternalServerError(message)

	}

	var auctionEntities []auction_entity.Auction

	for _, auctionEntityMongo := range auctionEntitiesMongo {
		auctionEntities = append(auctionEntities, auction_entity.Auction{
			Id:          auctionEntityMongo.Id,
			ProductName: auctionEntityMongo.ProductName,
			Description: auctionEntityMongo.Description,
			Condition:   auctionEntityMongo.Condition,
			Status:      auctionEntityMongo.Status,
			Category:    auctionEntityMongo.Category,
			Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
		})
	}

	return auctionEntities, nil
}
