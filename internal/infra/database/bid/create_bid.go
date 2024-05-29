package bid

import (
	"context"
	"fmt"
	"sync"

	"github.com/felipemagrassi/auction-go/configuration/logger"
	"github.com/felipemagrassi/auction-go/internal/entity/auction_entity"
	"github.com/felipemagrassi/auction-go/internal/entity/bid_entity"
	"github.com/felipemagrassi/auction-go/internal/infra/database/auction"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	Id        string  `bson:"id"`
	UserId    string  `bson:"user_id"`
	AuctionId string  `bson:"auction_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *auction.AuctionRepository
}

func NewBidRepository(database *mongo.Database, auctionRepository *auction.AuctionRepository) *BidRepository {
	return &BidRepository{
		Collection:        database.Collection("bids"),
		AuctionRepository: auctionRepository,
	}
}

func (bd *BidRepository) CreateBid(
	ctx context.Context,
	bidEntities []bid_entity.Bid,
) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, bid := range bidEntities {
		wg.Add(1)

		go func(bidValue bid_entity.Bid) {
			defer wg.Done()
			auctionEntity, err := bd.AuctionRepository.FindAuctionById(ctx, bidValue.AuctionId)
			if err != nil {
				message := fmt.Sprintf("Error trying to find auction with id = %s", bidValue.AuctionId)
				logger.Error(message, err)
				return
			}

			if auctionEntity.Status != auction_entity.Active {
				return
			}

			bidEntityMongo := &BidEntityMongo{
				Id:        bidValue.Id,
				UserId:    bidValue.UserId,
				AuctionId: bidValue.AuctionId,
				Amount:    bidValue.Amount,
				Timestamp: bidValue.Timestamp.Unix(),
			}

			if _, err := bd.Collection.InsertOne(ctx, bidEntityMongo); err != nil {
				message := fmt.Sprintf("Error trying to insert bid %v", bidValue)
				logger.Error(message, err)
				return
			}
		}(bid)

	}

	wg.Wait()
	return nil
}
