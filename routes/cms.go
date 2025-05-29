package routes

import (
	"encoding/json"
	"impel/cms-database/migrations"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func normalizeRows(rows []migrations.Entry) []map[string]interface{} {
	data, err := json.Marshal(rows)

	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("serialized.json", data, 0666); err != nil {
		log.Fatal(err)
	}

	groupedData := []map[string]interface{}{}

	for _, row := range rows {
		productID := row.Product.ID
		productName := row.Product.ProductName
		groupID := row.Group.ID
		groupName := row.Group.GroupName
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
				"Groups":      []map[string]interface{}{},
			}
			groupedData = append(groupedData, productEntry)
		}

		// Find or create the group entry within the product
		groups := productEntry["Groups"].([]map[string]interface{})
		var groupEntry map[string]interface{}
		for _, group := range groups {
			if group["GroupID"] == groupID {
				groupEntry = group
				break
			}
		}

		if groupEntry == nil {
			groupEntry = map[string]interface{}{
				"GroupID":   groupID,
				"GroupName": groupName,
				"Topics":    []map[string]interface{}{},
			}
			productEntry["Groups"] = append(groups, groupEntry)
		}

		// Find or create the topic entry within the group
		topics := groupEntry["Topics"].([]map[string]interface{})
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
			groupEntry["Topics"] = append(topics, topicEntry)
		}

		// Append the QA to the topic
		topicEntry["QAs"] = append(topicEntry["QAs"].([]map[string]interface{}), map[string]interface{}{
			"EntryID":  row.ID,
			"Question": row.Question.Question,
			"Answer":   row.Answer.Answer,
		})
	}

	data, err = json.Marshal(groupedData)

	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("normalized.json", data, 0666); err != nil {
		log.Fatal(err)
	}

	return groupedData
}

func GetCmsHandler(c *fiber.Ctx, db *gorm.DB) error {
	dealerId := c.Params("dealerId")

	// get rows for the dealer with joins
	var rows []migrations.Entry
	db.Preload("Dealer").
		Preload("Product").
		Preload("Group").
		Preload("Topic").
		Preload("Question").
		Preload("Answer").
		Order("ID desc").
		Where("dealer_id = ?", dealerId).
		Find(&rows)

	groupedData := normalizeRows(rows)

	return c.JSON(groupedData)
}
