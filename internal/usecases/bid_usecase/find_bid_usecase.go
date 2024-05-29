package bid_usecase

import (
	"context"

	"github.com/felipemagrassi/auction-go/internal/internal_error"
)

func (bu *BidUseCase) FindBidByAuctionId(
	ctx context.Context,
	auctionId string,
) ([]BidOutputDTO, *internal_error.InternalError) {
	bidEntities, err := bu.BidRepository.FindBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	bidOutputs := make([]BidOutputDTO, 0)
	for _, bidEntity := range bidEntities {
		bidOutputs = append(bidOutputs, BidOutputDTO{
			Id:        bidEntity.Id,
			AuctionId: bidEntity.AuctionId,
			UserId:    bidEntity.UserId,
			Amount:    bidEntity.Amount,
			Timestamp: bidEntity.Timestamp,
		})
	}

	return bidOutputs, nil
}
