package copy

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CopyTopicToSegmentRequestBody struct {
	DealerID     uint `json:"DealerID"`
	ProductID    uint `json:"ProductID"`
	SegmentID    uint `json:"SegmentID"`
	TopicID      uint `json:"TopicID"`
	NewSegmentID uint `json:"NewSegmentID"`
}

func CopyTopicToSegment(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CopyTopicToSegmentRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid JSON body",
		})
	}

	// Validate required fields
	if body.DealerID == 0 ||
		body.ProductID == 0 ||
		body.SegmentID == 0 ||
		body.TopicID == 0 ||
		body.NewSegmentID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing or invalid fields",
		})
	}

	// minimal security
	if body.SegmentID == body.NewSegmentID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "trying to copy to same segment",
		})
	}

	// get every QA record in this topic for the dealer
	var topicRecords []migrations.Entry
	var topicRecord = migrations.Entry{
		DealerID:  body.DealerID,
		ProductID: body.ProductID,
		SegmentID: body.SegmentID,
		TopicID:   body.TopicID,
	}
	db.Where(&topicRecord).Find(&topicRecords)

	// go over all of those records and copy them with the new segment id
	for _, record := range topicRecords {
		newRow := migrations.Entry{
			DealerID:   record.DealerID,
			ProductID:  record.ProductID,
			SegmentID:  body.NewSegmentID,
			TopicID:    record.TopicID,
			QuestionID: record.QuestionID,
			AnswerID:   record.AnswerID,
		}
		db.Create(&newRow)
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
