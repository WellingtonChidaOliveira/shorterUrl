package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/shortner/models"
	"github.com/wellingtonchida/shortner/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type Handlers struct {
	service service.ShorterService
}

func NewHandler(svc service.ShorterService) *Handlers {
	return &Handlers{
		service: svc,
	}
}

func (h *Handlers) GetAllShorter(c *fiber.Ctx) error {
	shorterList, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}
	return c.Status(http.StatusAccepted).JSON(shorterList)
}

func (h *Handlers) GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	shorter, err := h.service.GetById(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}
	return c.Status(http.StatusOK).JSON(shorter)
}

type requestCreate struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}


func (h *Handlers) Create(c *fiber.Ctx) error {
	var sh requestCreate
	if err := c.BodyParser(&sh); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	send := &models.Shorter{
		ID:        primitive.NewObjectID(),
		Shorter:   sh.Name,
		OriginUrl: sh.Uri,
		MostUse:   0,
		Status:    1,
		CreatedAt: time.Now(),
	}

	oid, err := h.service.Create(c.Context(), *send)
	if err != nil {

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": oid})
}
type requestUpdate struct {
	Name string `json:"name,omitempty"`
	Uri  string `json:"uri,omitempty"`
	Count  int `json:"use,omitempty"`
}
func (h *Handlers) Update(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.SendString("Need ID to get this Shorter")
	}

	 var sh requestUpdate
	if err := c.BodyParser(&sh); err != nil {
		return c.Status(http.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
	}

	send := &models.Shorter{
		Shorter: sh.Name,
		OriginUrl: sh.Uri,
		MostUse: sh.Count,
	}

	err := h.service.Update(c.Context(), id, *send)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusAccepted)
}

func (h *Handlers) Delete(c *fiber.Ctx) error {
	id := c.Query("id")

	err := h.service.Delete(c.Context(), id)
	if err != nil {
		return nil
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handlers) Inactivate(c *fiber.Ctx) error {
	id := c.Query("id")

	err := h.service.Inactivate(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handlers)Activate(c *fiber.Ctx) error {
	id := c.Query("id")

	err := h.service.Activate(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}




