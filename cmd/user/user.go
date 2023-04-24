package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adesupraptolaia/user_login/config"
	"github.com/adesupraptolaia/user_login/db"
	_ "github.com/adesupraptolaia/user_login/docs"
	user_profile_controller "github.com/adesupraptolaia/user_login/internal/controller/user_profile"
	"github.com/adesupraptolaia/user_login/internal/entity"
	"github.com/adesupraptolaia/user_login/internal/repo"
	"github.com/adesupraptolaia/user_login/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:8000
// @BasePath /
func Run() {
	cfg := config.Config

	db, err := db.NewDatabase("./db/migrations/user")
	if err != nil {
		log.Panicf("error when init database, err: %s", err.Error())
	}
	seedData(db)

	userProfileRepo := repo.NewUserProfile(db)
	authRepo := repo.NewAuthRepo()
	usecase := usecase.NewUserProfile(userProfileRepo, authRepo)

	publicHandler := user_profile_controller.NewUserProfileHandler(usecase)

	c := echo.New()

	c.Use(middleware.Logger())
	c.Use(middleware.Recover())

	c.GET("/", healthCheck)
	c.GET("/user/:user_ksuid", publicHandler.GetUser)
	c.POST("/user/create", publicHandler.CreateUser)
	c.POST("/user/:user_ksuid/update", publicHandler.UpdateUser)
	c.DELETE("/user/:user_ksuid", publicHandler.DeleteUser)

	c.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		if err := c.Start(fmt.Sprintf(":%d", cfg.UserServer.Port)); err != nil {
			log.Fatalf("Failed to start server, err: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully shut down both servers and any active connections
	if err := c.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down server, err: %s", err.Error())
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
	db.Table("user_profiles").Create(entity.UserProfile{
		UserKsuid:   adminKsuid,
		Name:        "Admin",
		DateOfBirth: time.Now().Format("2006-01-02"),
		Address:     "Perawang",
	})

	// create user
	db.Table("user_profiles").Create(entity.UserProfile{
		UserKsuid:   userKsuid,
		Name:        "User",
		DateOfBirth: "2019-01-01",
		Address:     "Malang",
	})
}
