package server

import (
	"encoding/json"
	"log"
	"new-go-project/cmd/service"

	"github.com/gofiber/fiber/v3"
)

func (s *server) handleCreateUser(c fiber.Ctx) error {
	ctx := c.Context()

	// parse request
	var user service.User
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		log.Println("failed to unmarshal request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	// set initial status to inactive
	user.Status = service.UserStatusInactive

	// create user
	err = s.userService.CreateUser(ctx, &user)
	if err != nil {
		log.Println("failed to create user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	// return response
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (s *server) handleGetUsers(c fiber.Ctx) error {
	ctx := c.Context()

	users, err := s.userService.GetUsers(ctx)
	if err != nil {
		log.Println("failed to get users", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get users",
		})
	}

	return c.JSON(users)
}

func (s *server) handleGetUserById(c fiber.Ctx) error {
	userId := c.Params("id")
	user, err := s.userService.GetUserById(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

func (s *server) handleActivateUser(c fiber.Ctx) error {
	userId := c.Params("id")

	user, err := s.userService.ActivateUser(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

func (s *server) handleDeleteUser(c fiber.Ctx) error {
	userId := c.Params("id")

	err := s.userService.DeleteUser(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User successfully deleted",
		"user_id": userId,
	})
}

type updateUserNameRequest struct {
	FirstName  *string `json:"first_name,omitempty"`
	MiddleName *string `json:"middle_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
}

func (s *server) handleUpdateUserName(c fiber.Ctx) error {
	userId := c.Params("id")

	var req updateUserNameRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validasi minimal satu field diisi
	if req.FirstName == nil && req.MiddleName == nil && req.LastName == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "at least one field (first_name, middle_name, or last_name) must be provided",
		})
	}

	// Validasi field yang diisi tidak boleh kosong
	if req.FirstName != nil && *req.FirstName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "first_name cannot be empty",
		})
	}

	if req.LastName != nil && *req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "last_name cannot be empty",
		})
	}

	user, err := s.userService.UpdateUserName(c.Context(), userId, req.FirstName, req.MiddleName, req.LastName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile successfully updated",
		"user":    user,
	})
}
