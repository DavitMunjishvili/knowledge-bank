package migrations

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/gorm"
)

type DefaultQA []DefaultQAElement

type DefaultQAElement struct {
	Product ProductName    `json:"product"`
	Topics  []DefaultTopic `json:"topics"`
}

type DefaultTopic struct {
	Label     string            `json:"label"`
	Questions []DefaultQuestion `json:"questions"`
}

type DefaultQuestion struct {
	Answer   string `json:"answer"`
	Question string `json:"question"`
}

func CreateDefaultQAs(db gorm.DB, dealerId uint, groupId uint) (ok bool) {
	ok = false
	var dealerRecord Dealer
	var groupRecord Group

	db.First(&dealerRecord, dealerId)
	db.First(&groupRecord, groupId)

	content, err := os.ReadFile("./migrations/default_qa.json")
	if err != nil {
		log.Fatal("Error during ReadFile", err)
		panic("Error when opening default_qa.json")
	}

	var payload DefaultQA
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal()", err)
		panic("Error when unmarshalling default_qa content")
	}

	for _, product := range payload {
		var productRecord = Product{ProductName: product.Product}
		db.First(&productRecord)

		for _, topic := range product.Topics {
			var topicRecord = Topic{
				TopicName: topic.Label,
			}
			db.Where(&topicRecord).FirstOrCreate(&topicRecord)

			for _, qa := range topic.Questions {
				var questionRecord = Question{
					Question: qa.Question,
				}
				db.Where(&questionRecord).FirstOrCreate(&questionRecord)

				var answerRecord = Answer{
					Answer: qa.Answer,
				}
				db.Where(&answerRecord).FirstOrCreate(&answerRecord)

				entry := Entry{
					Dealer:   dealerRecord,
					Group:    groupRecord,
					Product:  productRecord,
					Topic:    topicRecord,
					Question: questionRecord,
					Answer:   answerRecord,
				}
				db.Create(&entry)
			}
		}
	}
	ok = true
	return
}
