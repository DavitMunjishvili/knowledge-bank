package routes

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateTopicRequestBody struct {
	DealerID  uint   `json:"DealerID"`
	ProductID uint   `json:"ProductID"`
	SegmentID uint   `json:"SegmentID"`
	TopicName string `json:"TopicName"`
}

func PostTopicHandler(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CreateTopicRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid JSON body",
		})
	}

	// Validate required fields
	if body.DealerID == 0 || body.ProductID == 0 || body.SegmentID == 0 || body.TopicName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing or invalid fields",
		})
	}

	// find topic or create new one
	var topicRecord = migrations.Topic{TopicName: body.TopicName}
	db.Where(topicRecord).
		Attrs(&migrations.Topic{Custom: true}).
		FirstOrCreate(&topicRecord)

	// create new topic
	entry := migrations.Entry{
		DealerID:  body.DealerID,
		SegmentID: body.SegmentID,
		ProductID: body.ProductID,
		Topic:     topicRecord,
	}
	db.Create(&entry)

	return c.JSON(fiber.Map{
		"status": "success",
	})

}
