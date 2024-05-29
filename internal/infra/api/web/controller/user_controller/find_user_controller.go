package user_controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/felipemagrassi/auction-go/configuration/rest_err"
	"github.com/felipemagrassi/auction-go/internal/usecases/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userUsecase user_usecase.UserUseCaseInterface
}

func NewUserController(userUsecase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

func (u *UserController) FindUserById(c *gin.Context) {
	userId := c.Param("userId")
	if err := uuid.Validate(userId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid Fields", rest_err.Causes{
			Field:   "userId",
			Message: fmt.Sprintf("Invalid UUID value: %s", userId),
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	userData, err := u.userUsecase.FindUserById(context.Background(), userId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, userData)
}
