package auction_entity

import (
	"context"
	"time"

	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"github.com/google/uuid"
)

type AuctionRepositoryInterface interface {
	CreateAuction(context.Context, *Auction) *internal_error.InternalError
	FindAuctionById(context.Context, string) (*Auction, *internal_error.InternalError)
	FindAuctions(context.Context, AuctionStatus, string, string) ([]Auction, *internal_error.InternalError)
}

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

type (
	ProductCondition int
	AuctionStatus    int
)

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota
	Used
	Refurbished
)

func CreateAuction(
	productName, category, description string,
	condition ProductCondition,
) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (au *Auction) Validate() *internal_error.InternalError {
	if len(au.ProductName) <= 1 ||
		len(au.Category) <= 1 {
		return internal_error.NewBadRequestError("Invalid auction data")
	}

	return nil
}
