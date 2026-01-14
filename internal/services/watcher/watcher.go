package watcher

import (
	"context"
	"fmt"

	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/rabbitmq"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/domain"
)

type watcherService struct {
	fsClient *firestore.Client
	mqClient *rabbitmq.Client
}

func NewService(fsClient *firestore.Client, mqClient *rabbitmq.Client) *watcherService {
	// TODO: create new watcher service
	return &watcherService{
		fsClient: fsClient,
		mqClient: mqClient,
	}
}

func (s *watcherService) Start(ctx context.Context) {
	// TODO: listen to firestore collection for changes & call pushToQueue() if a new document isnt already in the queue
	snapIterator := s.fsClient.Client.Collection(s.fsClient.Collection).Snapshots(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			docs, err := snapIterator.Query.Where("status", "==", "PENDING").Documents(ctx).GetAll()
			if err != nil {
				fmt.Errorf("error getting documents: %d", err)
			}

			for _, val := range docs {
				data := val.Data()
			}
		}
	}
}

func pushToQueue(migrationJob domain.MigrationJob, mqClient rabbitmq.Client) {
	// TODO: push job to rabbitMQ queue & call updateJobToInProgress()
}

func updateJobToInProgress(fsClient firestore.Client, docId uint32) {
	// TODO: update document 'status' field to "IN_PROGRESS"
}
