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

// listens to firestore collection for changes & pushes to queue for each job sent through the channel
func (s *watcherService) Start(ctx context.Context) {
	log.Println("[WATCHER] Listening for pending jobs..")
	jobs, errors := s.fsClient.ListenForPendingJobs(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[WATCHER] Context done.")
			return
		case err := <-errors:
			log.Printf("[WATCHER] Error: %v", err)
		case job := <-jobs:
			log.Printf("[WATCHER] Found new job: %s", job.ID)
			s.queueNewJob(job, ctx)
		}
	}
}

// pushes job to rabbitMQ & updates job status to "QUEUED"
func (s *watcherService) queueNewJob(migrationJob domain.MigrationJob, ctx context.Context) {
	log.Printf("[WATCHER] Queuing new job: %s", migrationJob.ID)

	err := s.mqClient.PublishJob(migrationJob)

	if err != nil {
		log.Printf("[WATCHER] Error publishing job: %v", err)
		return
	}
	log.Printf("[WATCHER] Job %s published to queue.", migrationJob.ID)

	err = s.fsClient.UpdateJobStatus(ctx, migrationJob, domain.StatusQueued)
	if err != nil {
		log.Printf("[WATCHER] Error updating job status: %v", err)
		return
	}

	log.Printf("[WATCHER] Job %s status updated to %s.", migrationJob.ID, domain.StatusQueued)
}
