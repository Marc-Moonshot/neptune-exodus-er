package firestore

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	gfs "cloud.google.com/go/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/config"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/domain"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	Client     *gfs.Client
	Collection string
}

// wrapper to create and return a new firestore client
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	client, err := gfs.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.ApplicationCredentials))
	if err != nil {
		log.Println("Error initializing firestore client.")
		log.Fatal(err)
	}

	fsClient := Client{
		Client:     client,
		Collection: cfg.CollectionName,
	}
	return &fsClient, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

// listens for migration jobs with status of "PENDING" and return them
func (c *Client) ListenForPendingJobs(ctx context.Context) (<-chan domain.MigrationJob, <-chan error) {
	jobsChan := make(chan domain.MigrationJob)
	errChan := make(chan error)

	go func() {
		defer close(jobsChan)
		defer close(errChan)

		query := c.Client.Collection(c.Collection).Where("status", "==", domain.StatusPending).OrderBy("createdBy", gfs.Asc)

		snapIterator := query.Snapshots(ctx)
		defer snapIterator.Stop()

		firstSnapshot := true

		for {
			snap, err := snapIterator.Next()

			if err == iterator.Done {
				break
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errChan <- err:
					// broken stream, idk what to do here
					return
				}
			}

			if firstSnapshot {
				for {
					docSnap, err := snap.Documents.Next()
					if err == iterator.Done {
						break
					}
					if err != nil {
						errChan <- err
						break
					}

					var job domain.MigrationJob
					if err := docSnap.DataTo(&job); err != nil {
						errChan <- err
						continue
					}

					select {
					case <-ctx.Done():
						return
					case jobsChan <- job:
					}
				}
				firstSnapshot = false
				continue
			}

			for _, change := range snap.Changes {
				if change.Kind == gfs.DocumentAdded || change.Kind == gfs.DocumentModified {
					var job domain.MigrationJob
					if err := change.Doc.DataTo(&job); err != nil {
						errChan <- err
						continue
					}
					select {
					case <-ctx.Done():
						return
					case jobsChan <- job:
					}
				}
			}
		}
	}()
	return jobsChan, errChan
}

// updates job status to "IN_PROGRESS" in firestore collection
func (c *Client) UpdateJobStatus(ctx context.Context, job domain.MigrationJob) error {
	docRef := c.Client.Collection("data_migrations").Doc(job.ID)
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "status", Value: domain.StatusRunning},
	})

	if err != nil {
		return fmt.Errorf("Failed to update job status: %w", err)
	}
	log.Printf("Job %s status updated to %s", job.ID, domain.StatusRunning)
	return nil
}
