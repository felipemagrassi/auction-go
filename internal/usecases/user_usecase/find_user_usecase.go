package user_usecase

import (
	"context"

	"github.com/felipemagrassi/auction-go/internal/entity/user_entity"
	"github.com/felipemagrassi/auction-go/internal/internal_error"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserById(context.Context, string) (*UserOutputDTO, *internal_error.InternalError)
}

func (u *UserUseCase) FindUserById(
	ctx context.Context,
	userId string,
) (*UserOutputDTO, *internal_error.InternalError) {
	userEntity, err := u.UserRepository.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}, nil
}
