package copy

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CopyGroupToProductRequestBody struct {
	DealerID     uint `json:"DealerID"`
	ProductID    uint `json:"ProductID"`
	GroupID      uint `json:"GroupID"`
	NewProductID uint `json:"NewProductID"`
}

func CopyGroupToProduct(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CopyGroupToProductRequestBody)
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
		body.NewProductID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing or invalid fields",
		})
	}

	// minimal security
	if body.ProductID == body.NewProductID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "trying to copy to same group",
		})
	}

	// get every QA record in this group for the dealer
	var groups []migrations.Entry
	var group = migrations.Entry{
		DealerID:  body.DealerID,
		ProductID: body.ProductID,
		GroupID:   body.GroupID,
	}
	db.Where(&group).Find(&groups)

	// go over all of those records and copy them with the new product id
	for _, record := range groups {
		newRow := migrations.Entry{
			DealerID:   record.DealerID,
			ProductID:  body.NewProductID,
			GroupID:    record.GroupID,
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
