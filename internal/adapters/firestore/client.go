package firestore

import (
	"context"
	"log"

	gfs "cloud.google.com/go/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/config"
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

// func (c *Client) ListenForPendingJobs() {
// 	// TODO: listen for migration jobs with status of "PENDING" and return them
// }

// func (c *Client) UpdateJobStatus() {
// 	// TODO: update job status to "IN_PROGRESS" in firestore collection
// }
