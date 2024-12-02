package routes

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateDealerRequestBody struct {
	DealerName   string                 `json:"DealerName"`
	SalesforceID string                 `json:"SalesforceID"`
	Metadata     map[string]interface{} `json:"Metadata"`
	SegmentID    uint                   `json:"SegmentID"`
}

func PostDealerHandler(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CreateDealerRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid JSON body",
		})
	}

	if body.DealerName == "" || body.SalesforceID == "" || body.SegmentID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing or invalid fields",
		})
	}

	// create dealer
	dealer := migrations.Dealer{
		DealerName:   body.DealerName,
		SalesforceID: body.SalesforceID,
	}
	db.Create(&dealer)

	// create default QAs for the dealer
	// TODO: need to figure out SegmentID
	ok := migrations.CreateDefaultQAs(*db, dealer.ID, body.SegmentID)

	if ok {
		return c.JSON(fiber.Map{
			"status": "success",
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status": "fail",
		"error":  "idk what but something failed",
	})
}
