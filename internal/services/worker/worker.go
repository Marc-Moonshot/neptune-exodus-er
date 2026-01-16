package worker

import (
	"context"
	"log"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/rabbitmq"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/domain"
)

type workerService struct {
	mq rabbitmq.Client
}

func newService(rabbitmq *rabbitmq.Client) *workerService {
	return &workerService{
		mq: *rabbitmq,
	}
}

// Starts migration function depending on job
func (w *workerService) Start(ctx context.Context) {
	jobs, errs := w.mq.ConsumeJob(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("context done. worker shutting down.")
		case err := <-errs:
			if err != nil {
				log.Printf("consumer error: %v", err)
			}
		case job := <-jobs:
			w.handleJob(job)
		}
	}
}

func (w *workerService) handleJob(job domain.MigrationJob) {
	log.Printf("processing job %s", job.ID)
}
