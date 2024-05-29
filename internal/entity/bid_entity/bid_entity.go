package bid_entity

import (
	"context"
	"time"

	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"github.com/google/uuid"
)

type BidEntityRepositoryInterface interface {
	CreateBid(context.Context, []Bid) *internal_error.InternalError
	FindBidByAuctionId(context.Context, string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByAuctionId(context.Context, string) (*Bid, *internal_error.InternalError)
}

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

func CreateBid(userId, auctionId string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		UserId:    userId,
		AuctionId: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) Validate() *internal_error.InternalError {
	if b.Amount <= 0 {
		return internal_error.NewBadRequestError("amount must be greater than 0")
	}

	if err := uuid.Validate(b.UserId); err != nil {
		return internal_error.NewBadRequestError("invalid user id")
	}

	if err := uuid.Validate(b.AuctionId); err != nil {
		return internal_error.NewBadRequestError("invalid auction id")
	}

	return nil
}
