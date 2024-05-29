package auction_usecase

import (
	"context"
	"time"

	"github.com/felipemagrassi/auction-go/internal/entity/auction_entity"
	"github.com/felipemagrassi/auction-go/internal/entity/bid_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"github.com/felipemagrassi/auction-go/internal/usecases/bid_usecase"
)

type AuctionInputDTO struct {
	ProductName string                          `json:"product_name" binding:"required,min=1"`
	Category    string                          `json:"category" binding:"required,min=2"`
	Description string                          `json:"description" binding:"required,min=10"`
	Condition   auction_entity.ProductCondition `json:"condition"`
}

type AuctionOutputDTO struct {
	Id          string                          `json:"id"`
	ProductName string                          `json:"product_name"`
	Category    string                          `json:"category"`
	Description string                          `json:"description"`
	Condition   auction_entity.ProductCondition `json:"condition"`
	Status      auction_entity.AuctionStatus    `json:"status"`
	Timestamp   time.Time                       `json:"timestamp" time_format:"2006-01-02T15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

type AuctionUseCase struct {
	AuctionRepository auction_entity.AuctionRepositoryInterface
	BidRepository     bid_entity.BidEntityRepositoryInterface
}

type AuctionUseCaseInterface interface {
	CreateAuction(context.Context, AuctionInputDTO) *internal_error.InternalError
	FindAuctionById(context.Context, string) (*AuctionOutputDTO, error)
	FindAuctions(context.Context, auction_entity.AuctionStatus, string, string) ([]AuctionOutputDTO, error)
}

func (au *AuctionUseCase) CreateAuction(
	ctx context.Context,
	auctionInput AuctionInputDTO,
) *internal_error.InternalError {
	auction, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition),
	)
	if err != nil {
		return err
	}

	if err := au.AuctionRepository.CreateAuction(ctx, auction); err != nil {
		return err
	}

	return nil
}
