package middleware

import (
	"fiber-ngulik/pkg/httperror"
	"fiber-ngulik/pkg/sdk/jwtrsa"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func VerifyJWTRSA(publicKey string) MiddlewareFunc {
	verifyPublicKey, err := jwtrsa.GetPublicKey(publicKey)

	if err != nil {
		log.Default().Printf("%s", err.Error())
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{JWTAlg: "RS256", Key: verifyPublicKey},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": httperror.UnauthorizedErrorMessage,
			})
		},
	})
}
