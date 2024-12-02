package copy

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CopySegmentToProductRequestBody struct {
	DealerID     uint `json:"DealerID"`
	ProductID    uint `json:"ProductID"`
	SegmentID    uint `json:"SegmentID"`
	NewProductID uint `json:"NewProductID"`
}

func CopySegmentToProduct(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CopySegmentToProductRequestBody)
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
			"error":  "trying to copy to same segment",
		})
	}

	// get every QA record in this segment for the dealer
	var segments []migrations.Entry
	var segment = migrations.Entry{
		DealerID:  body.DealerID,
		ProductID: body.ProductID,
		SegmentID: body.SegmentID,
	}
	db.Where(&segment).Find(&segments)

	// go over all of those records and copy them with the new product id
	for _, record := range segments {
		newRow := migrations.Entry{
			DealerID:   record.DealerID,
			ProductID:  body.NewProductID,
			SegmentID:  record.SegmentID,
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
