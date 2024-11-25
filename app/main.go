package main

import (
	"fiber-ngulik/config"
	"fiber-ngulik/pkg/constant"
	"fiber-ngulik/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	userDomain "fiber-ngulik/modules/user/domain"
	userHandler "fiber-ngulik/modules/user/handlers"
	userRepository "fiber-ngulik/modules/user/repositories"
	userUsecase "fiber-ngulik/modules/user/usecases"

	authDomain "fiber-ngulik/modules/auth/domain"
	authHandler "fiber-ngulik/modules/auth/handlers"
	authUsecase "fiber-ngulik/modules/auth/usecases"

	databases "fiber-ngulik/pkg/databases"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type repositories struct {
	userRepository userDomain.UserRepository
}

type usecase struct {
	userUsecase userDomain.UserUsecase
	authUsecase authDomain.AuthUsecase
}

type packages struct {
	repositories repositories
	usecase      usecase
}

var pkg packages

func setPackages() {
	// repository
	pkg.repositories.userRepository = userRepository.NewUserRepository(databases.DBConnect.Connection)

	// usecase
	pkg.usecase.userUsecase = userUsecase.NewUserUsecase(pkg.repositories.userRepository)
	pkg.usecase.authUsecase = authUsecase.NewAuthUsecase(pkg.repositories.userRepository)

}

func setHttp(f *fiber.App) {
	f.Get("/"+constant.API_VERSION_1+"/health-check", func(c *fiber.Ctx) error {
		log.Default().Println("main", "This service is running properly")
		return utils.Response(nil, "This service is running properly", 200, c)
	})

	userHandler.NewUserHandler(f, pkg.usecase.userUsecase)

	authHandler.NewAuthHandler(f, pkg.usecase.authUsecase)
}

func main() {
	path, _ := os.Getwd()
	utils.LogDefault(path)

	databases.InitConnection(config.Config().PostgreDSN())

	app := fiber.New(fiber.Config{
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(logger.New(logger.Config{
		Next:       logger.ConfigDefault.Next,
		Format:     `[ROUTE] ${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}` + "\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(recover.New())
	setPackages()
	setHttp(app)

	app.Use(cors.New(cors.ConfigDefault))

	if err := app.Listen(":3000"); err != nil && err != http.ErrServerClosed {
		log.Default().Println("main", fmt.Sprintf("Could not listen on %s: %v\n", ":3000", err))
	}
}
