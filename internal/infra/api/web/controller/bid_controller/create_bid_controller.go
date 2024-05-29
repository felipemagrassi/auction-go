package bid_controller

import (
	"context"
	"net/http"

	"github.com/felipemagrassi/auction-go/configuration/rest_err"
	"github.com/felipemagrassi/auction-go/internal/infra/api/web/validation"
	"github.com/felipemagrassi/auction-go/internal/usecases/bid_usecase"
	"github.com/gin-gonic/gin"
)

type bidController struct {
	bidUseCase bid_usecase.BidUseCaseInterface
}

func NewBidController(bidUseCase bid_usecase.BidUseCaseInterface) *bidController {
	return &bidController{
		bidUseCase: bidUseCase,
	}
}

func (u *bidController) CreateBid(c *gin.Context) {
	var bidInputDTO bid_usecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.bidUseCase.CreateBid(context.Background(), bidInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
