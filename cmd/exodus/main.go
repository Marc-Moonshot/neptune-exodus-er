package main

import (
	"context"
	"log"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/rabbitmq"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/config"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/services/watcher"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Neptune ExodusER.")

	ctx := context.Background()

	// Adapter initialization
	var fsClient *firestore.Client

	fsClient, err = firestore.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init Firestore: %v", err)
	}
	defer fsClient.Close()

	log.Println("Initialized firestore.")

	var mqClient *rabbitmq.Client

	mqClient, err = rabbitmq.NewClient(cfg)
	if err != nil {
		log.Println("Failed to init rabbitMQ.")
		log.Fatalln(err)
	}

	defer mqClient.Close()
	log.Println("Initialized rabbitMQ.")

	// Run watcher in a goroutine, worker in main thread
	go runWatcher(ctx, fsClient, mqClient)
	<-ctx.Done()
}

func runWatcher(ctx context.Context, fs *firestore.Client, mq *rabbitmq.Client) {
	svc := watcher.NewService(fs, mq)
	svc.Start(ctx)
}

// func runWorker(fs *firestore.Client, mq *rabbitmq.Client) {
// 	// svc := worker.NewService(fs, mq)
// 	// svc.Start()
// }
