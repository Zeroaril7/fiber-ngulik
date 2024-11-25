package handlers

import (
	"fiber-ngulik/modules/auth/domain"
	"fiber-ngulik/modules/auth/models"
	"fiber-ngulik/pkg/httperror"
	"fiber-ngulik/pkg/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authUsecase domain.AuthUsecase
}

// Login implements AuthHandler.
func (h *authHandler) Login(c *fiber.Ctx) error {
	authRequest := new(models.LoginAuth)

	if err := c.BodyParser(authRequest); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	validate := validator.New()

	if err := validate.Struct(authRequest); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	result := <-h.authUsecase.AuthWithPassword(c.Context(), *authRequest)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Login success", http.StatusOK, c)
}

func NewAuthHandler(f *fiber.App, authUsecase domain.AuthUsecase) AuthHandler {
	handler := &authHandler{
		authUsecase: authUsecase,
	}

	group := f.Group("/auth")
	group.Post("/login", handler.Login)

	return handler
}
