package rabbitmq

import (
	"log"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewClient(cfg *config.Config) (*Client, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("Error creating amqp client.")
		log.Fatalln(err)
	}

	ch, chanErr := conn.Channel()
	if chanErr != nil {
		log.Printf("Error creating amqp channel.")
		log.Fatalln(chanErr)
	}

	q, qErr := ch.QueueDeclare(
		"migration-jobs",
		true,
		false,
		true,
		false,
		nil,
	)
	if qErr != nil {
		log.Printf("Error creating amqp queue.")
		log.Fatalln(qErr)
	}

	mqClient := Client{
		conn:    conn,
		channel: ch,
		queue:   q,
	}

	return &mqClient, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PublishJob() {
	// TODO: push new job to rabbitMQ queue
}

func (c *Client) ConsumeJob() {
	// TODO: consume job from rabbitMQ queueu
}
