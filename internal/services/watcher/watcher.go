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

func (s *watcherService) Start(ctx context.Context) {
	// TODO: listen to firestore collection for changes & call pushToQueue() if a new document isnt already in the queue
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
			s.processNewJob(job, ctx)
		}
	}
}

func pushToQueue(migrationJob domain.MigrationJob, mqClient rabbitmq.Client) {
	// TODO: push job to rabbitMQ queue & call updateJobToInProgress()
}

func updateJobToInProgress(fsClient firestore.Client, docId uint32) {
	// TODO: update document 'status' field to "IN_PROGRESS"
}
