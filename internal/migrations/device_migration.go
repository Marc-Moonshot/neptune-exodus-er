package migrations

import (
	"context"
	"fmt"
	"log"
	"strings"

	gfs "cloud.google.com/go/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/adapters/firestore"
	"github.com/Marc-Moonshot/neptune-exodus-er/internal/domain"
	"google.golang.org/api/iterator"
)

func DeviceMigration(job domain.MigrationJob, fs firestore.Client, ctx context.Context) error {
	var fromDevice = job.FromDevice
	var toDevice = job.ToDevice

	fromDev, toDev, err := getDeviceMetadata(ctx, fs, fromDevice, toDevice)
	if err != nil {
		return err
	}

	// Query
	var collection = fs.Client.Collection("device_data")
	query := collection.Where("deviceId", "==", fromDevice).OrderBy("day", gfs.Asc)
	if !job.MigrateAllData {
		query = query.Where("day", ">=", job.FromDate).Where("day", "<=", job.ToDate)
	}

	bw := fs.Client.BulkWriter(ctx)
	defer bw.End()

	docIterator := query.Documents(ctx)
	targetCollection := fs.Client.Collection("device_data")

	// Stream documents one by one
	for {
		docSnap, err := docIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("iterator error: %v", err)
		}

		var d domain.Data
		if err := docSnap.DataTo(&d); err != nil {
			log.Printf("Failed to parse document %s: %v", docSnap.Ref.ID, err)
			continue
		}

		// Perform conversion on the single document in memory
		if job.ConvertData {
			performUnitConversion(&d, fromDev, toDev)
		}

		parts := strings.Split(docSnap.Ref.ID, "_")
		if len(parts) < 2 {
			continue
		}
		newDocName := fmt.Sprintf("%d_%s", toDevice, parts[1])
		docRef := targetCollection.Doc(newDocName)

		// Prepare Data Map
		data := map[string]any{
			"deviceId":  toDevice,
			"day":       d.Day,
			"pressure":  d.Pressure,
			"flow_rate": d.Flow_rate,
			"net_flow":  d.Net_flow,
			"errors":    d.Errors,
		}

		// 5. Schedule Write
		if job.OverrideData {
			bw.Set(docRef, data, gfs.MergeAll)
		} else {
			bw.Create(docRef, data)
		}
	}

	return nil
}

// in memory implementation
// func DeviceMigration(job domain.MigrationJob, fs firestore.Client, ctx context.Context) error {
// 	var fromDevice = job.FromDevice
// 	var toDevice = job.ToDevice

// 	// gather all documents in memory
// 	// TODO: make collection name dynamic via a new field in each migration
// 	var collection = fs.Client.Collection("device_data")

// 	docs := collection.Where("deviceId", "==", fromDevice).OrderBy("day", gfs.Asc)
// 	if !job.MigrateAllData {
// 		docs = docs.Where("day", ">=", job.FromDate).Where("day", "<=", job.ToDate)
// 	}

// 	docIterator := docs.Documents(ctx)

// 	var allDocs []domain.Data
// 	for {
// 		doc, err := docIterator.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			log.Printf("Failed to iterate: %v", err)
// 			continue
// 		}
// 		var d domain.Data
// 		if err := doc.DataTo(&d); err != nil {
// 			log.Printf("Failed to parse document %s: %v", doc.Ref.ID, err)
// 			continue
// 		}
// 		d.Id = doc.Ref.ID
// 		allDocs = append(allDocs, d)
// 	}

// 	fromQuery := fs.Client.Collection("io_devices").Where("device_number", "==", fromDevice).Limit(1).Documents(ctx)
// 	fromSnap, err := fromQuery.Next()
// 	if err == iterator.Done {
// 		return fmt.Errorf("no source device found for device_number %d", fromDevice)
// 	}
// 	if err != nil {
// 		return fmt.Errorf("failed to query source device: %v", err)
// 	}

// 	toQuery := fs.Client.Collection("io_devices").Where("device_number", "==", toDevice).Limit(1).Documents(ctx)
// 	toSnap, err := toQuery.Next()
// 	if err == iterator.Done {
// 		return fmt.Errorf("no target device found for device_number %d", toDevice)
// 	}
// 	if err != nil {
// 		return fmt.Errorf("failed to query target device: %v", err)
// 	}

// 	var fromDev, toDev domain.Device
// 	if err := fromSnap.DataTo(&fromDev); err != nil {
// 		return fmt.Errorf("failed to parse source device: %v", err)
// 	}
// 	if err := toSnap.DataTo(&toDev); err != nil {
// 		return fmt.Errorf("failed to parse target device: %v", err)
// 	}

// 	if job.ConvertData {
// 		for i, doc := range allDocs {
// 			for j, dp := range doc.Pressure {
// 				allDocs[i].Pressure[j].Value = int64(convertValue(
// 					float64(dp.Value),
// 					fromDev.Units[0],
// 					toDev.Units[0],
// 				))
// 			}
// 			for j, dp := range doc.Flow_rate {
// 				allDocs[i].Flow_rate[j].Value = int64(convertValue(
// 					float64(dp.Value),
// 					fromDev.Units[0],
// 					toDev.Units[0],
// 				))
// 			}
// 			for j, dp := range doc.Net_flow {
// 				allDocs[i].Net_flow[j].Value = int64(convertValue(
// 					float64(dp.Value),
// 					fromDev.Units[0],
// 					toDev.Units[0],
// 				))
// 			}
// 		}
// 	}

// 	bw := fs.Client.BulkWriter(ctx)
// 	defer bw.End()

// 	targetCollection := fs.Client.Collection("device_data")

// 	for _, doc := range allDocs {
// 		parts := strings.Split(doc.Id, "_")
// 		if len(parts) < 2 {
// 			log.Printf("Skipping doc with unexpected ID format: %s", doc.Id)
// 			continue
// 		}
// 		datePart := parts[1]
// 		newDocName := fmt.Sprintf("%d_%s", toDevice, datePart)

// 		docRef := targetCollection.Doc(newDocName)

// 		data := map[string]any{
// 			"deviceId":  toDevice,
// 			"day":       doc.Day,
// 			"pressure":  doc.Pressure,
// 			"flow_rate": doc.Flow_rate,
// 			"net_flow":  doc.Net_flow,
// 			"errors":    doc.Errors,
// 		}

// 		if job.OverrideData {
// 			_, err := bw.Set(docRef, data, gfs.MergeAll)
// 			if err != nil {
// 				log.Printf("Bulkwriter set error for %s: %v", doc.Id, err)
// 			}
// 		} else {
// 			docSnap, err := docRef.Get(ctx)
// 			if err != nil {
// 				if !docSnap.Exists() {
// 					_, err := bw.Create(docRef, data)
// 					if err != nil {
// 						log.Printf("Bulkwriter set error for %s: %v", doc.Id, err)
// 					}
// 				} else {
// 					log.Printf("Failed to check existence for %s: %v", doc.Id, err)
// 				}
// 			} else {
// 				log.Printf("Skipping existing document: %s", doc.Id)
// 			}
// 		}
// 	}
// 	return nil
// }

func getDeviceMetadata(ctx context.Context, fs firestore.Client, fromID, toID int64) (domain.Device, domain.Device, error) {
	var fromDev, toDev domain.Device

	fetch := func(id int64) (domain.Device, error) {
		query := fs.Client.Collection("io_devices").Where("device_number", "==", id).Limit(1).Documents(ctx)
		defer query.Stop()

		snap, err := query.Next()
		if err == iterator.Done {
			return domain.Device{}, fmt.Errorf("device %d not found", id)
		}
		if err != nil {
			return domain.Device{}, err
		}

		var dev domain.Device
		if err := snap.DataTo(&dev); err != nil {
			return domain.Device{}, err
		}
		return dev, nil
	}

	// Fetch both in parallel
	var err1, err2 error
	done := make(chan bool)

	go func() { fromDev, err1 = fetch(fromID); done <- true }()
	go func() { toDev, err2 = fetch(toID); done <- true }()

	// Wait for both goroutines
	<-done
	<-done

	if err1 != nil {
		return fromDev, toDev, fmt.Errorf("source: %w", err1)
	}
	if err2 != nil {
		return fromDev, toDev, fmt.Errorf("target: %w", err2)
	}

	return fromDev, toDev, nil
}

func performUnitConversion(d *domain.Data, fromDev, toDev domain.Device) {

	for i := range d.Pressure {
		d.Pressure[i].Value = int64(convertValue(
			float64(d.Pressure[i].Value),
			fromDev.Units[0],
			toDev.Units[0],
		))
	}

	for i := range d.Flow_rate {
		d.Flow_rate[i].Value = int64(convertValue(
			float64(d.Flow_rate[i].Value),
			fromDev.Units[0],
			toDev.Units[0],
		))
	}

	for i := range d.Net_flow {
		d.Net_flow[i].Value = int64(convertValue(
			float64(d.Net_flow[i].Value),
			fromDev.Units[0], toDev.Units[0],
		))
	}
}

func convertValue(value float64, fromUnit string, toUnit string) float64 {
	if fromUnit == toUnit {
		return value
	}

	fromUnit = strings.ToLower(fromUnit)
	toUnit = strings.ToLower(toUnit)

	switch fromUnit {

	//  PRESSURE
	case "psi":
		switch toUnit {
		case "bar":
			return value * 0.0689476
		case "kpa":
			return value * 6.89476
		}

	case "bar":
		switch toUnit {
		case "psi":
			return value * 14.5038
		case "kpa":
			return value * 100
		}

	case "kpa":
		switch toUnit {
		case "psi":
			return value * 0.145038
		case "bar":
			return value * 0.01
		}

	// FLOW RATE
	case "m3hr", "cu.m/hr", "m³/h":
		switch toUnit {
		case "lmin", "l/min":
			return value * 16.6667
		}

	case "lmin", "l/min":
		switch toUnit {
		case "m3hr", "cu.m/hr", "m³/h":
			return value / 16.6667
		}

	//  FLOW NET
	case "m3", "cu.m":
		switch toUnit {
		case "l", "liters":
			return value * 1000
		}

	case "l", "liters":
		switch toUnit {
		case "m3", "cu.m":
			return value / 1000
		}
	}

	log.Printf("Unrecognized device units: %s to %s", fromUnit, toUnit)
	return value
}
