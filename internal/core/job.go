// TODO
// [ ] - Remove job status and type(topic) as topic is determined the queue it is in and then the status is encoded in the tcp protocol

package core

type JobStatus string

// Status of the job

const (
	StatusPending    JobStatus = "PENDING"    // Not yet picked up by any worker
	StatusProcessing JobStatus = "PROCESSING" // Picked up by a worker and the worker started processing the job
	StatusCompleted  JobStatus = "COMPLETED"  // The worker completed the processing of the job
	StatusFailed     JobStatus = "FAILED"     // The worker failed to process the job
)

// Job Struct
// ID - id of the job, auto incremented.
// Type - type of the job, which determines the queue in which it waits. Like "image-resize", "video-encoding", "email-sender".
// Status - Current status of the job.
// Payload - Payload of the Job as an array of bytes (can be converted to json)

type Job struct {
	ID      uint64
	Topic   string
	Status  JobStatus
	Payload []byte
}