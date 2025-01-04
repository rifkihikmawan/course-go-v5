package server

import (
	"new-go-project/cmd/internal/postgres"
	"new-go-project/cmd/service"

	"github.com/gofiber/fiber/v3"
)

type Server interface {
	Start()
}

type server struct {
	app *fiber.App
	db  *postgres.Client

	userService service.UserService
}

func NewServer() Server {
	app := fiber.New()

	db, err := postgres.NewClient()
	if err != nil {
		panic(err)
	}

	userStore := postgres.NewUserStore(db)
	userService := service.NewUserService(userStore)

	return &server{
		app: app,
		db:  db,

		userService: userService,
	}
}

func (s *server) Start() {
	s.app.Get("/hello-world", s.handleHelloWorld)
	s.app.Get("/hello/:name", s.handleHello)

	// health check
	s.app.Get("/health", s.handleHealth)

	// users
	s.app.Post("/users", s.handleCreateUser)
	s.app.Get("/users", s.handleGetUsers)
	s.app.Get("/users/:id", s.handleGetUserById)
	s.app.Post("/users/:id/activate", s.handleActivateUser)
	s.app.Delete("/users/:id", s.handleDeleteUser)
	s.app.Patch("/users/profile/:id", s.handleUpdateUserName)

	s.app.Listen(":8080")
}
