package main

import (
	"context"
	"log"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Neptune ExodusER.")
	log.Printf("Queue URL: %s", cfg.RabbitMQURL)
	log.Printf("Collection name: %s", cfg.CollectionName)

	ctx := context.Background()

	// Adapter initialization
	var fsClient *firestore.Client

	fsClient, err = firestore.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init Firestore: %v", err)
	}
	defer fsClient.Close()

	log.Println("Initialized firestore.")

	// 2. figure out how to connect to rabbitMQ
	// var mqClient *rabbitmq.Client

	// Run watcher in a goroutine, worker in main thread

}

// func runWatcher(ctx context.Context, fs *firestore.Client, mq *rabbitmq.Client) {
// svc := watcher.NewService(fs, mq)
// svc.Start(ctx)
// }

// func runWorker(fs *firestore.Client, mq *rabbitmq.Client) {
// 	// svc := worker.NewService(fs, mq)
// 	// svc.Start()
// }
