package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/wellingtonchida/shortner/handler"
	"github.com/wellingtonchida/shortner/repository"
	"github.com/wellingtonchida/shortner/service"
)

func main() {
	ctx := context.Background()
	repo := repository.NewDataBase(ctx)
	app := fiber.New()
	svc := service.Service(repo)

	h := handler.NewHandler(svc)
	startRoutes(app, *h)
	app.Listen(":3000")
}

func startRoutes(app fiber.Router, h handler.Handlers) {
	app.Get("/", h.GetAllShorter)
	app.Get("/:id", h.GetById)
	app.Post("/create", h.Create)
	app.Put("/update", h.Update)
	app.Delete("/", h.Delete)
	app.Patch("/", h.Inactivate)
}
