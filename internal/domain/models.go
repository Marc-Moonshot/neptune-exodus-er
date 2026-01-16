package domain

type JobStatus string

const (
	StatusPending JobStatus = "PENDING"
	StatusQueued  JobStatus = "QUEUED"
	StatusRunning JobStatus = "RUNNING"
	StatusDone    JobStatus = "DONE"
	StatusFailed  JobStatus = "FAILED"
)

// migration task document from Firestore
type MigrationJob struct {
	ID             string    `firestore:"-" json:"id"`
	FromDevice     int64     `firestore:"fromDevice" json:"from_device"`
	ToDevice       int64     `firestore:"toDevice" json:"to_device"`
	Status         JobStatus `firestore:"status" json:"status"`
	FromDate       int64     `firestore:"fromDate,omitempty" json:"from_date"`
	ToDate         int64     `firestore:"toDate,omitempty" json:"to_date"`
	ConvertData    bool      `firestore:"convertData" json:"convert_data"`
	MigrateAllData bool      `firestore:"migrateAllData" json:"migrate_all_data"`
	OverrideData   bool      `firestore:"createdBy" json:"override_data"`
	CreatedAt      int64     `firestore:"createdAt" json:"created_at"`
	CreatedBy      string    `firestore:"createdBy" json:"created_by"`
	Error          string    `firestore:"error,omitempty" json:"error,omitempty"`
}
