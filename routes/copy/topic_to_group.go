package copy

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CopyTopicToGroupRequestBody struct {
	DealerID   uint `json:"DealerID"`
	ProductID  uint `json:"ProductID"`
	GroupID    uint `json:"GroupID"`
	TopicID    uint `json:"TopicID"`
	NewGroupID uint `json:"NewGroupID"`
}

func CopyTopicToGroup(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CopyTopicToGroupRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid JSON body",
		})
	}

	// Validate required fields
	if body.DealerID == 0 ||
		body.ProductID == 0 ||
		body.GroupID == 0 ||
		body.TopicID == 0 ||
		body.NewGroupID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing or invalid fields",
		})
	}

	// minimal security
	if body.GroupID == body.NewGroupID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "trying to copy to same group",
		})
	}

	// get every QA record in this topic for the dealer
	var topicRecords []migrations.Entry
	var topicRecord = migrations.Entry{
		DealerID:  body.DealerID,
		ProductID: body.ProductID,
		GroupID:   body.GroupID,
		TopicID:   body.TopicID,
	}
	db.Where(&topicRecord).Find(&topicRecords)

	// go over all of those records and copy them with the new group id
	for _, record := range topicRecords {
		newRow := migrations.Entry{
			DealerID:   record.DealerID,
			ProductID:  record.ProductID,
			GroupID:    body.NewGroupID,
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
