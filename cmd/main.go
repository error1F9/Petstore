package main

import (
	_ "Petstore/docs"
	"Petstore/models/pet/controller"
	"Petstore/models/pet/repository"
	"Petstore/models/pet/service"
	storeController "Petstore/models/store/controller"
	store "Petstore/models/store/repository"
	storeService "Petstore/models/store/service"
	userController "Petstore/models/user/controller"
	user "Petstore/models/user/repository"
	userService "Petstore/models/user/service"
	"Petstore/pkg/config"
	"Petstore/pkg/dbase"
	"Petstore/pkg/migrator"
	"Petstore/pkg/token"
	"Petstore/responder"
	"Petstore/router"
	"context"
	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var decoder = godecoder.NewDecoder(jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	DisallowUnknownFields:  true,
})

// @title			My Petstore
// @version		1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	support@example.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8080
// @BasePath		/

// @tag.name		pet
// @tag.description	Access to Petstore orders

// @tag.name		store
// @tag.description	Everything about your Pets

// @tag.name		user
// @tag.description	Operations about user
func main() {
	logger, _ := zap.NewProduction()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := dbase.NewPostgersDB(cfg)
	if err != nil {
		log.Fatalf("Error initializing PostgersDB: %v", err)
	}

	if err = migrator.RunMigrations(db); err != nil {
		log.Fatalf("Error initializing PostgersDB: %v", err)
	}

	tokenService := token.NewJWTTokenService("fuioasjdf89as%#(_Qj890g423wj90-yj4wegh09j4we0")
	resp := responder.NewResponder(decoder, logger)

	petRepo := repository.NewPetRepository(db)
	petService := service.NewPetService(petRepo)
	petController := controller.NewPetController(petService, resp, decoder)

	storeRepo := store.NewStoreRepository(db)
	newStoreService := storeService.NewOrderService(storeRepo)
	newStoreController := storeController.NewStoreController(newStoreService, resp, decoder)

	userRepo := user.NewUserRepository(db)
	newUserService := userService.NewUserService(userRepo, *tokenService)
	newUserController := userController.NewUserController(newUserService, resp, decoder)

	initHandlers := router.InitRoutes(petController, newStoreController, newUserController, tokenService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: initHandlers,
	}
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped")
}
