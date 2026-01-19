package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/config"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type ConsumedJob struct {
	Job domain.MigrationJob
	Msg amqp.Delivery
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
		false,
		false,
		amqp.Table{
			"x-message-deduplication": true,
		},
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

// pushes new job to rabbitMQ queue
func (c *Client) PublishJob(job domain.MigrationJob) error {

	body, err := json.Marshal(job)

	if err != nil {
		return fmt.Errorf("failed to unmarshall job object: %w", err)
	}
	c.channel.Publish("", c.queue.Name, false, false, amqp.Publishing{
		MessageId:    job.ID,
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
	return nil
}

// Consumes jobs from rabbitMQ queue and returns them one by one
func (c *Client) ConsumeJob(ctx context.Context) (chan ConsumedJob, chan error) {
	jobs := make(chan ConsumedJob)
	errs := make(chan error)

	msgs, err := c.channel.Consume(c.queue.Name, "", false, false, false, false, nil)

	if err != nil {
		errs <- fmt.Errorf("failed to start consumer: %w", err)
		close(jobs)
		close(errs)
		return jobs, errs
	}
	log.Printf("Consuming jobs from queue %s", c.queue.Name)

	go func() {
		defer close(jobs)
		defer close(errs)

		for {
			select {
			case <-ctx.Done():
				log.Println("context done, stopping.")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("message channel closed.")
					return
				}

				var job domain.MigrationJob

				if err := json.Unmarshal(msg.Body, &job); err != nil {
					log.Printf("failed to parse body: %v", err)
					msg.Nack(false, false)
					continue
				}
				log.Printf("Emitting job: %s", job.ID)

				jobs <- ConsumedJob{Job: job, Msg: msg}
				// msg.Ack(false)
			}
		}
	}()
	return jobs, errs
}

func (c *Client) Acknowledge(msg amqp.Delivery) error {
	msg.Ack(false)
	return nil
}
