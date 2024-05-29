package bid_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/felipemagrassi/auction-go/configuration/logger"
	"github.com/felipemagrassi/auction-go/internal/entity/bid_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
)

var bidBatch []bid_entity.Bid

type BidOutputDTO struct {
	Id        string    `json:"id"`
	AuctionId string    `json:"auction_id"`
	UserId    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidUseCaseInterface interface {
	CreateBid(context.Context, BidInputDTO) *internal_error.InternalError
	FindBidByAuctionId(context.Context, string) ([]BidOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(context.Context, string) (*BidOutputDTO, *internal_error.InternalError)
}

type BidUseCase struct {
	BidRepository bid_entity.BidEntityRepositoryInterface

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration

	bidChannel chan bid_entity.Bid
}

func NewBidUseCase(
	bidRepository bid_entity.BidEntityRepositoryInterface,
) *BidUseCase {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()
	bidUseCase := &BidUseCase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
		timer:               time.NewTimer(maxSizeInterval),
	}

	bidUseCase.triggerCreateBidRoutine(context.Background())

	return bidUseCase
}

func (bu *BidUseCase) triggerCreateBidRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)

		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("error trying to process bid batch list", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bidEntity)
				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}

					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("error trying to process bid batch list", err)
				}
				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)
			}
		}
	}()
}

func (bu *BidUseCase) CreateBid(
	ctx context.Context,
	bidInputDTO BidInputDTO,
) *internal_error.InternalError {
	bidEntity, err := bid_entity.CreateBid(
		bidInputDTO.UserId,
		bidInputDTO.AuctionId,
		bidInputDTO.Amount,
	)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bidEntity

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	maxBatchSize := os.Getenv("MAX_BATCH_SIZE")
	value, err := strconv.Atoi(maxBatchSize)
	if err != nil {
		return 5
	}

	return value
}
