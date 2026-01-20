package migrations

import (
	"context"
	"log"
	"time"

	gfs "cloud.google.com/go/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/domain"
	"google.golang.org/api/iterator"
)

func DeviceMigration(job domain.MigrationJob, fs firestore.Client, ctx context.Context) error {
	var fromDevice = job.FromDevice
	var toDevice = job.ToDevice
	var MMDDYYYY = time.Now().Format("01-02-2006")
	var newDocName = string(fromDevice) + "_" + MMDDYYYY

	// gather all documents in memory
	// TODO: make collection name dynamic via a new field in each migration
	var collection = fs.Client.Collection("device_data")

	docs := collection.Where("deviceId", "==", fromDevice).OrderBy("day", gfs.Asc)
	if !job.MigrateAllData {
		docs = docs.Where("day", ">=", job.FromDate).Where("day", "<=", job.ToDate)
	}

	docIterator := docs.Documents(ctx)

	var allDocs []domain.Data
	for {
		doc, err := docIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		var d domain.Data
		if err := doc.DataTo(&d); err != nil {
			log.Printf("Failed to parse document %s: %v", doc.Ref.ID, err)
			continue
		}

		allDocs = append(allDocs, d)
	}

	if job.ConvertData {
		// convert datapoints
		for _, doc := range allDocs {
			datapoints := doc.Pressure
			for _, datapoint := range datapoints {

			}
		}
	}

	if job.OverrideData {

	}

	return nil
}
