package pet

import (
	"Petstore/models/pet/controller"
	"Petstore/models/pet/entity"
	"Petstore/models/pet/repository"
	"Petstore/models/pet/service"
	storeController "Petstore/models/store/controller"
	store "Petstore/models/store/repository"
	storeService "Petstore/models/store/service"
	userController "Petstore/models/user/controller"
	user "Petstore/models/user/repository"
	userService "Petstore/models/user/service"
	"Petstore/pkg/dbase"
	"Petstore/pkg/migrator"
	"Petstore/pkg/token"
	"Petstore/responder"
	"Petstore/router"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var decoder = godecoder.NewDecoder(jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	DisallowUnknownFields:  true,
})

func setupMockSeverWithDB(t *testing.T) *chi.Mux {
	logger, _ := zap.NewProduction()

	db, err := dbase.NewSqLiteDB()
	if err != nil {
		log.Fatalf("Error initializing PostgersDB: %v", err)
	}

	if err = migrator.RunMigrationsNew(db); err != nil {
		log.Fatalf("Error initializing PostgersDB: %v", err)
	}

	tokenService := token.NewJWTTokenService("testkey123")
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

	return initHandlers
}

func TestAdd(t *testing.T) {
	tests := []struct {
		TestName string
		input    controller.AddPetRequest
		expected string
	}{
		{
			TestName: "Successful adding pet 1",
			input: controller.AddPetRequest{
				Name: "TestPet1",
				Category: entity.Category{
					Name: "TestCategory1",
				},
				Status: "available",
			},
		},
		{
			TestName: "Successful adding pet 2 (same category)",
			input: controller.AddPetRequest{
				Name: "TestPet2",
				Category: entity.Category{
					Name: "TestCategory1",
				},
				Status: "available",
			},
		},
		{
			TestName: "Successful adding pet 3 (diff category)",
			input: controller.AddPetRequest{
				Name: "TestPet2",
				Category: entity.Category{
					Name: "TestCategory2",
				},
				Status: "available",
			},
		},
	}

	r := setupMockSeverWithDB(t)

	server := httptest.NewServer(r)
	defer server.Close()

	client := http.Client{}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			data, _ := json.Marshal(test.input)
			fmt.Println(data)
			body := bytes.NewReader(data)
			resp, err := client.Post(server.URL+"/pet", "application/json", body)
			defer resp.Body.Close()
			if err != nil {
				t.Fatalf("Error post /pet: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected 200 OK got %d", resp.StatusCode)
			}
		})
	}
}
