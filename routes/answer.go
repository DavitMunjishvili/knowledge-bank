package routes

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateAnswerRequestBody struct {
	EntryID uint   `json:"EntryID"`
	Answer  string `json:"Answer"`
	Update  bool   `json:"Update"`
}

func PostAnswerHandler(c *fiber.Ctx, db *gorm.DB) error {
	body := new(CreateAnswerRequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Invalid JSON body",
		})
	}

	// Validate required fields
	if body.EntryID == 0 || body.Answer == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"error":  "Missing or invalid fields",
		})
	}

	if body.Update {
		// ***
		// when user is trying to update this answer for
		// ***

		// get the row that we are trying to update
		var row migrations.Entry
		db.First(&row, body.EntryID)

		// find answer or create new one
		var newAnswer = migrations.Answer{
			Answer: body.Answer,
		}
		db.Where(&newAnswer).
			Attrs(&migrations.Answer{Custom: true}).
			FirstOrCreate(&newAnswer)

		// find every row for the dealer with this answer
		var entriesWithAnswer []migrations.Entry
		search := migrations.Entry{
			DealerID:   row.DealerID,
			QuestionID: row.QuestionID,
		}
		db.Where(&search).Find(&entriesWithAnswer)

		// go over those rows and give update answer
		for _, entry := range entriesWithAnswer {
			entry.Answer = newAnswer
			db.Save(&entry)
		}

		return c.JSON(fiber.Map{
			"status": "success",
		})

	} else {
		// ***
		// when user is trying to update answer only for this question
		// ***

		// find answer or create new one
		var answerRecord = migrations.Answer{
			Answer: body.Answer,
		}
		db.Where(&answerRecord).
			Attrs(&migrations.Answer{Custom: true}).
			FirstOrCreate(&answerRecord)

		// get row for this QA
		var row migrations.Entry
		db.First(&row, body.EntryID)

		// update answer
		row.Answer = answerRecord
		db.Save(&row)

		return c.JSON(fiber.Map{
			"status": "success",
		})
	}

}
