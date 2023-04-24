package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/adesupraptolaia/user_login/docs"
	"gorm.io/gorm"

	"github.com/adesupraptolaia/user_login/config"
	"github.com/adesupraptolaia/user_login/db"
	user_controller "github.com/adesupraptolaia/user_login/internal/controller/user"
	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/adesupraptolaia/user_login/internal/repo"
	"github.com/adesupraptolaia/user_login/internal/usecase"
	"github.com/adesupraptolaia/user_login/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger Auth App - Public
// @version 1.0
// @description This is a API documentation for Auth APP
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:9000
// @BasePath /
func Run() {
	cfg := config.Config

	db, err := db.NewDatabase("./db/migrations/auth")
	if err != nil {
		log.Panicf("error when init database, err: %s", err.Error())
	}

	seedData(db)

	// define repo, usecase, and userHandler
	repo := repo.NewUser(db)
	usecase := usecase.NewUser(repo)
	userHandler := user_controller.NewUserHandler(usecase)

	// Public Server
	publicServer := echo.New()
	publicServer.Use(middleware.Logger())
	publicServer.Use(middleware.Recover())

	publicServer.GET("/", healthCheck)
	publicServer.GET("/refresh", userHandler.RefreshToken)
	publicServer.POST("/login", userHandler.Login)

	publicServer.GET("/swagger/*", echoSwagger.WrapHandler)

	// Private Server
	privateServer := echo.New()
	privateServer.Use(middleware.Logger())
	privateServer.Use(middleware.Recover())

	privateServer.GET("/", healthCheck)
	privateServer.POST("/user/create", userHandler.CreateUser)
	privateServer.DELETE("/user/:ksuid", userHandler.DeleteUser)

	privateServer.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		if err := publicServer.Start(fmt.Sprintf(":%d", cfg.AuthServer.Port.Public)); err != nil {
			log.Fatalf("Failed to start public server, err: %s", err.Error())
		}
	}()

	go func() {
		if err := privateServer.Start(fmt.Sprintf(":%d", cfg.AuthServer.Port.Private)); err != nil {
			log.Fatalf("Failed to start private server, err: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shut down both servers and any active connections
	if err := publicServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down public server, err: %s", err.Error())
	}
	if err := privateServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down private server, err: %s", err.Error())
	}
	log.Println("Servers shut down successfully.")
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

func seedData(db *gorm.DB) {
	adminKsuid := "2OokWa2yDw7yi7o9RpsAl58xuoW"
	userKsuid := "2OokWdyzR17GBzVsF6auODTuSxz"

	// create admin
	db.Table("users").Create(entity.User{
		Ksuid:    adminKsuid,
		Username: "admin",
		Password: utils.HashPassword("admin"),
		Role:     entity.ADMIN,
	})

	// create user
	db.Table("users").Create(entity.User{
		Ksuid:    userKsuid,
		Username: "user",
		Password: utils.HashPassword("user"),
		Role:     entity.USER,
	})
}
