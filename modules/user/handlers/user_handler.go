package handlers

import (
	"fiber-ngulik/config"
	"fiber-ngulik/middleware"
	"fiber-ngulik/modules/user/domain"
	"fiber-ngulik/modules/user/models"
	"fiber-ngulik/pkg/httperror"
	"fiber-ngulik/pkg/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Add(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	GetByUsername(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(f *fiber.App, userUsecase domain.UserUsecase) UserHandler {
	handler := &userHandler{
		userUsecase: userUsecase,
	}

	const paramUsername = "/:username"

	group := f.Group("/user")
	group.Delete(paramUsername, handler.Delete)
	group.Get("", middleware.VerifyJWTRSA(config.Config().PublicKey), handler.Get)
	group.Get(paramUsername, middleware.VerifyBasicAuth(config.Config().BasicAuthUsername, config.Config().BasicAuthPassword), handler.GetByUsername)
	group.Post("", handler.Add)
	group.Put(paramUsername, handler.Update)

	return handler
}

// Add implements UserHandler.
func (h *userHandler) Add(c *fiber.Ctx) error {
	data := new(models.UserAdd)

	if err := c.BodyParser(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	validate := validator.New()

	if err := validate.Struct(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	password := utils.HashPassword(data.Password)
	data.Password = password

	expend := models.User{}
	expend = data.ToUser(expend)

	result := <-h.userUsecase.Add(c.Context(), expend)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Add user success", http.StatusOK, c)
}

// Delete implements UserHandler.
func (h *userHandler) Delete(c *fiber.Ctx) error {
	username := utils.ConvertString(c.Params("username"))

	result := <-h.userUsecase.Delete(c.Context(), username)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(nil, "Delete User success", http.StatusOK, c)
}

// Get implements UserHandler.
func (h *userHandler) Get(c *fiber.Ctx) error {
	filter := new(models.UserFilter)

	if err := c.QueryParser(filter); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	if !filter.DisablePagination {
		filter.SetDefault()
	}

	result := <-h.userUsecase.Get(c.Context(), *filter)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.ResponseWithPagination(result.Data, "Get user success", http.StatusOK, result.Total, filter.GetPaginationRequest(), c)
}

// GetByUsername implements UserHandler.
func (h *userHandler) GetByUsername(c *fiber.Ctx) error {
	username := utils.ConvertString(c.Params("username"))

	result := <-h.userUsecase.GetByUsername(c.Context(), username)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Get user success", http.StatusOK, c)
}

// Update implements UserHandler.
func (h *userHandler) Update(c *fiber.Ctx) error {
	username := utils.ConvertString(c.Params("username"))

	result := <-h.userUsecase.GetByUsername(c.Context(), username)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	expend := result.Data.(models.User)

	if expend == (models.User{}) {
		return utils.ResponseError(httperror.NotFound(httperror.NotFoundErrorMessage), c)
	}

	data := new(models.UserAdd)
	if err := c.BodyParser(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(httperror.BindErrorMessage), c)
	}

	validate := validator.New()

	if err := validate.Struct(data); err != nil {
		return utils.ResponseError(httperror.BadRequest(err.Error()), c)
	}

	if !utils.CheckPasswordHash(data.Password, expend.Password) {
		newPassword := utils.HashPassword(data.Password)
		data.Password = newPassword
	}

	expend = data.ToUser(expend)

	result = <-h.userUsecase.Update(c.Context(), expend)

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Update user success", http.StatusOK, c)
}
