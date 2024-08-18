package main

import (
	"loan-engine/config"
	"loan-engine/handler"
	"loan-engine/repository"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Get()
	app := fiber.New()

	loanRepo := repository.New(cfg)
	loanHandler := handler.New(loanRepo, cfg.State)

	app.Post("/propose", loanHandler.Propose)
	app.Post("/approve", loanHandler.Approve)
	app.Post("/invest", loanHandler.Invest)
	app.Post("/disburse", loanHandler.Disburse)
	app.Get("/all-loan", loanHandler.List)

	app.Listen(":3000")
}
