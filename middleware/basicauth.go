package middleware

import (
	"fiber-ngulik/pkg/httperror"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

type MiddlewareFunc func(c *fiber.Ctx) error

func VerifyBasicAuth(username, password string) MiddlewareFunc {
	return basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string) bool {
			return user == username && pass == password
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": httperror.UnauthorizedErrorMessage,
			})
		},
		Realm: "Forbidden",
	})
}
