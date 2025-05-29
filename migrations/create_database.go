package migrations

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Dealer struct {
	gorm.Model
	SalesforceID string
	DealerName   string
	Metadata     map[string]interface{} `gorm:"serializer:json"`
}

type ProductName string

const (
	Chat_AI       ProductName = "Chat AI"
	Car_Buying_AI ProductName = "Car Buying AI"
	Sales_AI      ProductName = "Sales AI"
	Service_AI    ProductName = "Service AI"
)

type Product struct {
	gorm.Model
	ProductName ProductName `gorm:"unique"`
}

type GroupName string

const (
	Car_Buying GroupName = "Car Buying"
	Sales      GroupName = "Sales"
	Service    GroupName = "Service"
	Upsell     GroupName = "Upsell"
)

type GroupType string

const (
	Segment GroupType = "Segment"
)

// this is the only table that should support bilingual content
type Group struct {
	gorm.Model
	GroupName GroupName `gorm:"unique"`
	GroupType GroupType `gorm:"default:Segment"`
}

type Topic struct {
	gorm.Model
	TopicName string `gorm:"unique"`
	Custom    bool   `gorm:"default:false"`
}

type Question struct {
	gorm.Model
	Question string
	Custom   bool `gorm:"default:false"`
}

type Answer struct {
	gorm.Model
	Answer string
	Custom bool `gorm:"default:false"`
}

type Locales string

const (
	EN Locales = "EN"
	ES Locales = "ES"
	FR Locales = "FR"
)

type Entry struct {
	gorm.Model
	Locale     Locales
	DealerID   uint
	GroupID    uint
	ProductID  uint
	TopicID    uint
	QuestionID *uint
	AnswerID   *uint

	Dealer   Dealer
	Group    Group
	Product  Product
	Topic    Topic
	Question Question
	Answer   Answer
}

func Initialize() *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=420"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Dealer{})
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Group{})
	db.AutoMigrate(&Topic{})
	db.AutoMigrate(&Question{})
	db.AutoMigrate(&Answer{})
	db.AutoMigrate(&Entry{})

	products := []*Product{
		{ProductName: Chat_AI},
		{ProductName: Car_Buying_AI},
		{ProductName: Sales_AI},
		{ProductName: Service_AI},
	}

	for _, product := range products {
		db.Where(product).FirstOrCreate(&product)
	}

	groups := []*Group{
		{GroupName: Car_Buying, GroupType: Segment},
		{GroupName: Sales, GroupType: Segment},
		{GroupName: Service, GroupType: Segment},
		{GroupName: Upsell, GroupType: Segment},
	}

	for _, group := range groups {
		db.Where(group).FirstOrCreate(&group)
	}

	defaultTopics := []*Topic{
		{TopicName: "Financing", Custom: false},
		{TopicName: "Trade-In", Custom: false},
		{TopicName: "Discounts, promotions", Custom: false},
		{TopicName: "Taxes and Fees", Custom: false},
		{TopicName: "Price Negotiation", Custom: false},
		{TopicName: "Hold Car / Deposit", Custom: false},
		{TopicName: "Test drive at home", Custom: false},
		{TopicName: "Shipping", Custom: false},
		{TopicName: "Warranty", Custom: false},
		{TopicName: "General information about dealership", Custom: false},
		{TopicName: "Leasing, Custom", Custom: false},
		{TopicName: "Selling the car", Custom: false},
		{TopicName: "Info about AI assistant", Custom: false},
		{TopicName: "Dealership policies", Custom: false},
		{TopicName: "Insuranc", Custom: false},
		{TopicName: "Service", Custom: false},
	}

	for _, topic := range defaultTopics {
		db.Where(topic).FirstOrCreate(&topic)
	}

	return db
}
