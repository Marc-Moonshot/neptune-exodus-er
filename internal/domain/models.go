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
	ID             string    `firestore:"-" json:"-"`
	DeviceType     string    `firestore:"deviceType" json:"device_type"`
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
	Type           string    `firestore:"type,omitempty" json:"type,omitempty"`
}

type dataPoint struct {
	Timestamp int64 `firestore:"timestamp"`
	Value     int64 `firestore:"value"`
}

type errorPoint struct {
	ErrorCode string `firestore:"errorCode"`
	Timestamp int64  `firestore:"timestamp"`
}

type Device struct {
	Id              string   `firestore:"id"`
	Device_number   int64    `firestore:"device_number"`
	Nickname        string   `firestore:"nickname"`
	Last_status     string   `firestore:"last_status"`
	Last_contact    int64    `firestore:"last_contact"`
	Date_registered int64    `firestore:"date_registered"`
	Rtu_assigned    int64    `firestore:"rtu_assigned"`
	Device_type     string   `firestore:"device_type"`
	Device_code     string   `firestore:"device_code"`
	Live            bool     `firestore:"live"`
	Units           []string `firestore:"units"`
}

type Data struct {
	Id        string       `firestore:"-"`
	Day       int64        `firestore:"day"`
	DeviceId  int64        `firestore:"deviceId"`
	Pressure  []dataPoint  `firestore:"pressure"`
	Flow_rate []dataPoint  `firestore:"flow_rate"`
	Net_flow  []dataPoint  `firestore:"net_flow"`
	Errors    []errorPoint `firestore:"errors"`
}
