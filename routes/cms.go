package routes

import (
	"impel/cms-database/migrations"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func normalizeRows(rows []migrations.Entry) []map[string]interface{} {
	groupedData := []map[string]interface{}{}

	for _, row := range rows {
		productID := row.Product.ID
		productName := row.Product.ProductName
		segmentID := row.Segment.ID
		segmentName := row.Segment.SegmentName
		topicID := row.Topic.ID
		topicName := row.Topic.TopicName

		// Find or create the product entry
		var productEntry map[string]interface{}
		for _, product := range groupedData {
			if product["ProductID"] == productID {
				productEntry = product
				break
			}
		}

		if productEntry == nil {
			productEntry = map[string]interface{}{
				"ProductID":   productID,
				"ProductName": productName,
				"Segments":    []map[string]interface{}{},
			}
			groupedData = append(groupedData, productEntry)
		}

		// Find or create the segment entry within the product
		segments := productEntry["Segments"].([]map[string]interface{})
		var segmentEntry map[string]interface{}
		for _, segment := range segments {
			if segment["SegmentID"] == segmentID {
				segmentEntry = segment
				break
			}
		}

		if segmentEntry == nil {
			segmentEntry = map[string]interface{}{
				"SegmentID":   segmentID,
				"SegmentName": segmentName,
				"Topics":      []map[string]interface{}{},
			}
			productEntry["Segments"] = append(segments, segmentEntry)
		}

		// Find or create the topic entry within the segment
		topics := segmentEntry["Topics"].([]map[string]interface{})
		var topicEntry map[string]interface{}
		for _, topic := range topics {
			if topic["TopicID"] == topicID {
				topicEntry = topic
				break
			}
		}

		if topicEntry == nil {
			topicEntry = map[string]interface{}{
				"TopicID":   topicID,
				"TopicName": topicName,
				"QAs":       []map[string]interface{}{},
			}
			segmentEntry["Topics"] = append(topics, topicEntry)
		}

		// Append the QA to the topic
		topicEntry["QAs"] = append(topicEntry["QAs"].([]map[string]interface{}), map[string]interface{}{
			"EntryID":  row.ID,
			"Question": row.Question.Question,
			"Answer":   row.Answer.Answer,
		})
	}

	return groupedData
}

func GetCmsHandler(c *fiber.Ctx, db *gorm.DB) error {
	dealerId := c.Params("dealerId")

	// get rows for the dealer with joins
	var rows []migrations.Entry
	db.Preload("Dealer").
		Preload("Product").
		Preload("Segment").
		Preload("Topic").
		Preload("Question").
		Preload("Answer").
		Order("ID desc").
		Where("dealer_id = ?", dealerId).
		Find(&rows)

	groupedData := normalizeRows(rows)

	return c.JSON(groupedData)
}
