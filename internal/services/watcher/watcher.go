package watcher

import (
	"context"
	"log"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/rabbitmq"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/domain"
)

type watcherService struct {
	fsClient *firestore.Client
	mqClient *rabbitmq.Client
}

func NewService(fsClient *firestore.Client, mqClient *rabbitmq.Client) *watcherService {
	return &watcherService{
		fsClient: fsClient,
		mqClient: mqClient,
	}
}

// listens to firestore collection for changes & call pushToQueue() if a new document isnt already in the queue
func (s *watcherService) Start(ctx context.Context) {
	log.Println("Listening for pending jobs..")
	jobs, errors := s.fsClient.ListenForPendingJobs(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("context done.")
			return
		case err := <-errors:
			log.Printf("Watcher error: %v", err)
		case job := <-jobs:
			log.Printf("Watcher found new job: %s", job.ID)
			s.queueNewJob(job, ctx)
		}
	}
}

// pushes job to rabbitMQ & updates job status to "QUEUED"
func (s *watcherService) queueNewJob(migrationJob domain.MigrationJob, ctx context.Context) {
	log.Printf("queuing new job: %s", migrationJob.ID)

	err := s.mqClient.PublishJob(migrationJob)

	if err != nil {
		log.Printf("Error publishing job: %v", err)
		return
	}
	log.Printf("job %s published to queue.", migrationJob.ID)

	err = s.fsClient.UpdateJobStatus(ctx, migrationJob, domain.StatusQueued)
	if err != nil {
		log.Printf("Error updating job status: %v", err)
		return
	}

	log.Printf("job %s status updated to %s.", migrationJob.ID, domain.StatusQueued)
}
