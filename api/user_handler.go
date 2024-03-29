package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/vihaan404/hotel-go/db"
	"github.com/vihaan404/hotel-go/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "user not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.UserParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	err := h.userStore.DeleteUser(c.Context(), userId)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"deleted": userId})
}

func (h *UserHandler) HandlerUpdateUser(c *fiber.Ctx) error {
	var params types.UpdatedUserParams
	userId := c.Params("id")
	err := c.BodyParser(&params)
	if err != nil {
		return err
	}
	err = h.userStore.UpdateUser(c.Context(), userId, params)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userId})
}
