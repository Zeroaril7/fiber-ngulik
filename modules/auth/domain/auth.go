package domain

import (
	"context"
	"fiber-ngulik/modules/auth/models"
	"fiber-ngulik/pkg/utils"
)

type AuthUsecase interface {
	AuthWithPassword(ctx context.Context, authReq models.LoginAuth) <-chan utils.Result
}
