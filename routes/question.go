package routes

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateQuestionRequestBody struct {
	DealerID  uint   `json:"DealerID"`
	ProductID uint   `json:"ProductID"`
	GroupID   uint   `json:"GroupID"`
	TopicID   uint   `json:"TopicID"`
	Question  string `json:"Question"`
}

func PostQuestionHandler(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CreateQuestionRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid JSON body",
		})
	}

	// Validate required fields
	if body.DealerID == 0 || body.ProductID == 0 || body.GroupID == 0 || body.TopicID == 0 || body.Question == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing or invalid fields",
		})
	}

	// get question or create new one
	var questionRecord = migrations.Question{
		Question: body.Question,
	}
	db.Where(&questionRecord).
		Attrs(&migrations.Question{Custom: true}).
		FirstOrCreate(&questionRecord)

	// insert new question
	entry := migrations.Entry{
		DealerID:  body.DealerID,
		GroupID:   body.GroupID,
		ProductID: body.ProductID,
		TopicID:   body.TopicID,
		Question:  questionRecord,
	}
	res := db.Create(&entry)

	if res.Error != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status": "fail",
			"error":  res.Error,
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}
