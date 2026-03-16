// package queue

// import (
// 	"testing"

// 	"github.com/Bharanidharan2006/cabbageMQ/internal/core"
// )

// func TestQueue(t *testing.T) {
// 	queueManager := NewQueueManager()

// 	if queueManager == nil {
// 		t.Fatalf("Queue Manager cannot be created!")
// 	}

// 	queueManager.CreateQueue(TopicConfig{
// 		Topic: "email",
// 		Size: 5,
// 	})

// 	err := queueManager.Push("email", &core.Job{
// 		ID: 5,
// 	})

// 	if err != nil {
// 		t.Errorf("Cannot push a job into the queue")
// 	}

// 	job, ok, err := queueManager.Read("email") 

// 	if err != nil && !ok {
// 		t.Errorf("Cannot push a job into the queue")
// 	}

// 	if job.ID != 5 {
// 		t.Errorf("Job is not retrieved properly")
// 	}

// 	job, ok, err = queueManager.Read("email")

// 	if err != nil {
// 		t.Errorf("Queue should be empty, but error returned")
// 	}

// 	if ok {
// 		t.Errorf("Queue should be empty but ok = true is returned")
// 	}


// }