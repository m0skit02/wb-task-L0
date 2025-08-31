package kafka

import (
	"context"
	"encoding/json"
	"log"
	"wb-task-L0/pkg/cache"
	"wb-task-L0/pkg/models"
	"wb-task-L0/pkg/repository"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader    *kafka.Reader
	orderRepo repository.Order // теперь интерфейс, а не конкретная структура
	cache     *cache.OrderCache
}

func NewConsumer(brokers []string, topic, groupID string, repo repository.Order, cache *cache.OrderCache) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{
		reader:    reader,
		orderRepo: repo,
		cache:     cache,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started...")
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		// десериализация сообщения в order
		var order models.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("invalid message, cannot unmarshal: %v", err)
			continue
		}

		// сохраняем в БД с транзакцией
		if err := c.orderRepo.Create(ctx, &order); err != nil {
			log.Printf("failed to save order in DB: %v", err)
			continue
		}

		// подтверждаем сообщение в Kafka
		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("failed to commit message: %v", err)
		}

		// кладем в кэш
		c.cache.Set(order)

		log.Printf("order saved: %s", order.OrderUID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
