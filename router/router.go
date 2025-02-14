package router

import (
	"Petstore/models/pet/controller"
	store "Petstore/models/store/controller"
	user "Petstore/models/user/controller"
	"Petstore/pkg/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func InitRoutes(c *controller.PetControl, s *store.StoreControl, u *user.UserController, t *token.JWTTokenService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/pet", func(r chi.Router) {
		r.Post("/", c.Add)
		r.Put("/", c.Update)
		r.Get("/findByStatus", c.FindByStatus)
		r.Get("/{id}", c.FindById)
		r.Post("/{id}", c.UpdateById)
		r.Delete("/{id}", c.Delete)
	})
	r.Route("/store", func(r chi.Router) {
		r.Get("/inventory", s.Inventory)
		r.Route("/order", func(r chi.Router) {
			r.Post("/", s.PlaceOrder)
			r.Get("/{id}", s.FindOrderById)
			r.Delete("/{id}", s.DeleteById)
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/login", u.Login)
		r.Post("/createWithArray", u.CreateUsersWithArray)
		r.Post("/", u.CreateUser)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(t.GetJWTAuth()))
			r.Use(jwtauth.Authenticator(t.GetJWTAuth()))

			r.Post("/logout", u.Logout)
			r.Get("/{username}", u.GetUser)
			r.Put("/{username}", u.UpdateUser)
			r.Delete("/{username}", u.DeleteUser)
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return r
}
