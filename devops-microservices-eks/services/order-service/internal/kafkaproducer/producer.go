package kafkaproducer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/devopsdemo/order-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

// Shape must match OrderCreatedEvent on the Notification (Spring Boot) side —
// this is the contract between the two services.
type OrderCreatedEvent struct {
	OrderID       string `json:"orderId"`
	CustomerEmail string `json:"customerEmail"`
}

func New(cfg *config.Config) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBootstrapServers),
		Topic:    cfg.OrderCreatedTopic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: writer}
}

func (p *Producer) PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	})
	if err != nil {
		// NOTE: intentionally not retried here. Interview topic to learn next:
		// the "outbox pattern" — right now the DB insert (order created) and
		// this publish are two separate operations, not one atomic unit. If
		// the process crashes between them, the order exists but no event
		// was published ("dual write" problem). The outbox pattern fixes this
		// by writing the event to an outbox table in the SAME DB transaction
		// as the order, then a separate relay process publishes it to Kafka.
		slog.Error("kafka_publish_failed", "orderId", event.OrderID, "error", err)
		return err
	}

	slog.Info("order_created_event_published", "orderId", event.OrderID)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
