package usecases

import (
	"context"
	"fiber-ngulik/modules/user/domain"
	"fiber-ngulik/modules/user/models"
	"fiber-ngulik/pkg/httperror"
	"fiber-ngulik/pkg/utils"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

// Add implements domain.UserUsecase.
func (u *userUsecase) Add(ctx context.Context, data models.User) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.userRepository.Add(ctx, data)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result.Username}
	}()

	return output
}

// Delete implements domain.UserUsecase.
func (u *userUsecase) Delete(ctx context.Context, username string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		err := u.userRepository.Delete(ctx, username)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{}
	}()

	return output
}

// Get implements domain.UserUsecase.
func (u *userUsecase) Get(ctx context.Context, filter models.UserFilter) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, total, err := u.userRepository.Get(ctx, filter)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
			return
		}

		output <- utils.Result{Data: result, Total: total}
	}()

	return output
}

// GetByUsername implements domain.UserUsecase.
func (u *userUsecase) GetByUsername(ctx context.Context, username string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.userRepository.GetByUsername(ctx, username)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

// Update implements domain.UserUsecase.
func (u *userUsecase) Update(ctx context.Context, data models.User) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		result, err := u.userRepository.Update(ctx, data)

		if err != nil {
			output <- utils.Result{Error: httperror.InternalServerError(err.Error())}
		}

		output <- utils.Result{Data: result}
	}()

	return output
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepository: userRepository}
}
