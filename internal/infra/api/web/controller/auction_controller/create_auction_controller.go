package auction_controller

import (
	"context"
	"net/http"

	"github.com/felipemagrassi/auction-go/configuration/rest_err"
	"github.com/felipemagrassi/auction-go/internal/infra/api/web/validation"
	"github.com/felipemagrassi/auction-go/internal/usecases/auction_usecase"
	"github.com/gin-gonic/gin"
)

type auctionController struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *auctionController {
	return &auctionController{
		auctionUseCase: auctionUseCase,
	}
}

func (u *auctionController) CreateAuction(c *gin.Context) {
	var auctionInputDTO auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.auctionUseCase.CreateAuction(context.Background(), auctionInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
