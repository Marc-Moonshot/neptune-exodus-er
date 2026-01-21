package worker

import (
	"context"
	"log"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/rabbitmq"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/migrations"
)

type workerService struct {
	mqClient rabbitmq.Client
	fsClient firestore.Client
}

func newService(rabbitmq *rabbitmq.Client, firestore *firestore.Client) *workerService {
	return &workerService{
		mqClient: *rabbitmq,
		fsClient: *firestore,
	}
}

// Starts migration function depending on job
func (w *workerService) Start(ctx context.Context) {
	jobs, errs := w.mqClient.ConsumeJob(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("context done. worker shutting down.")
		case err := <-errs:
			if err != nil {
				log.Printf("consumer error: %v", err)
			}
		case job := <-jobs:
			w.handleJob(job, ctx)
		}
	}
}

func (w *workerService) handleJob(consumedJob rabbitmq.ConsumedJob, ctx context.Context) {
	log.Printf("processing job %s", consumedJob.Job.ID)

	var err error = nil
	if consumedJob.Job.Type == "" || consumedJob.Job.Type == "device_migration" {
		err = migrations.DeviceMigration(consumedJob.Job, w.fsClient, ctx)
	}

	if err != nil {
		log.Printf("migration error: %v", err)
		consumedJob.Msg.Nack(false, true)
	}
	consumedJob.Msg.Ack(false)
}
