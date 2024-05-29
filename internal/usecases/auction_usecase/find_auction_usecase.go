package auction_usecase

import (
	"context"

	"github.com/felipemagrassi/auction-go/internal/entity/auction_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
	"github.com/felipemagrassi/auction-go/internal/usecases/bid_usecase"
)

func (au *AuctionUseCase) FindAuctionById(
	ctx context.Context,
	auctionId string,
) (*AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.AuctionRepository.FindAuctionById(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	return &AuctionOutputDTO{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auction_entity.ProductCondition(auctionEntity.Condition),
		Status:      auction_entity.AuctionStatus(auctionEntity.Status),
		Timestamp:   auctionEntity.Timestamp,
	}, nil
}

func (au *AuctionUseCase) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category, productName string,
) ([]AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntities, err := au.AuctionRepository.FindAuctions(ctx,
		auction_entity.AuctionStatus(status),
		category,
		productName,
	)
	if err != nil {
		return nil, err
	}

	auctionOutputs := make([]AuctionOutputDTO, 0)
	for _, auctionEntity := range auctionEntities {
		auctionOutputs = append(auctionOutputs, AuctionOutputDTO{
			Id:          auctionEntity.Id,
			ProductName: auctionEntity.ProductName,
			Category:    auctionEntity.Category,
			Description: auctionEntity.Description,
			Condition:   auction_entity.ProductCondition(auctionEntity.Condition),
			Status:      auction_entity.AuctionStatus(auctionEntity.Status),
			Timestamp:   auctionEntity.Timestamp,
		})
	}

	return auctionOutputs, nil
}

func (au *AuctionUseCase) FindWinningBidByAuctionId(
	ctx context.Context,
	auctionId string,
) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	auction, err := au.AuctionRepository.FindAuctionById(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	auctionOutput := AuctionOutputDTO{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction_entity.ProductCondition(auction.Condition),
		Status:      auction_entity.AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}

	bidEntity, err := au.BidRepository.FindWinningBidByAuctionId(ctx, auction.Id)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutput,
			Bid:     nil,
		}, nil
	}

	bidOutput := &bid_usecase.BidOutputDTO{
		Id:        bidEntity.Id,
		AuctionId: bidEntity.AuctionId,
		UserId:    bidEntity.UserId,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp,
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutput,
		Bid:     bidOutput,
	}, nil
}
