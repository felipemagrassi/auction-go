package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/felipemagrassi/auction-go/configuration/logger"
	"github.com/felipemagrassi/auction-go/internal/entity/bid_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (bd *BidRepository) FindBidByAuctionId(
	ctx context.Context,
	auctionId string,
) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auctionId": auctionId}

	cursor, err := bd.Collection.Find(ctx, filter)
	if err != nil {
		message := fmt.Sprintf("Error trying to find bid by auctionId = %s", auctionId)
		logger.Error(message, err)
		return nil, internal_error.NewInternalServerError(message)
	}

	var bidEntitiesMongo []BidEntityMongo
	if err := cursor.All(ctx, &bidEntitiesMongo); err != nil {
		message := fmt.Sprintf("Error trying to decode bid by auctionId = %s", auctionId)
		logger.Error(message, err)
		return nil, internal_error.NewInternalServerError(message)
	}

	var bidEntities []bid_entity.Bid
	for _, bidEntityMongo := range bidEntitiesMongo {
		bidEntities = append(bidEntities, bid_entity.Bid{
			Id:        bidEntityMongo.Id,
			UserId:    bidEntityMongo.UserId,
			AuctionId: bidEntityMongo.AuctionId,
			Amount:    bidEntityMongo.Amount,
			Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
		})
	}

	return bidEntities, nil
}

func (bd *BidRepository) FindWinningBidByAuctionId(
	ctx context.Context,
	auctionId string,
) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auctionId": auctionId}
	opts := options.FindOne().SetSort(bson.D{{"amount", -1}})

	var bidEntityMongo BidEntityMongo
	if err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		message := fmt.Sprintf("Error trying to find bid by auctionId = %s", auctionId)
		logger.Error(message, err)
		return nil, internal_error.NewInternalServerError(message)
	}

	return &bid_entity.Bid{
		Id:        bidEntityMongo.Id,
		UserId:    bidEntityMongo.UserId,
		AuctionId: bidEntityMongo.AuctionId,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
