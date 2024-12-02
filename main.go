package main

import (
	"impel/cms-database/migrations"
	"impel/cms-database/routes"
	"impel/cms-database/routes/copy"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db := migrations.Initialize()

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/cms/:dealerId", func(c *fiber.Ctx) error {
		return routes.GetCmsHandler(c, db)
	})

	app.Post("/dealer", func(c *fiber.Ctx) error {
		return routes.PostDealerHandler(c, db)
	})

	app.Post("/topic", func(c *fiber.Ctx) error {
		return routes.PostTopicHandler(c, db)
	})

	app.Post("/question", func(c *fiber.Ctx) error {
		return routes.PostQuestionHandler(c, db)
	})

	app.Post("/answer", func(c *fiber.Ctx) error {
		return routes.PostAnswerHandler(c, db)
	})

	app.Post("/copy-topic-to-segment", func(c *fiber.Ctx) error {
		return copy.CopyTopicToSegment(c, db)
	})

	app.Post("/copy-segment-to-product", func(c *fiber.Ctx) error {
		return copy.CopySegmentToProduct(c, db)
	})

	app.Listen(":3000")
}
