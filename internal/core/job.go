package core

type Status string

const ()

type Job struct {
	ID      int
	Type    string
	Status  Status
	Payload []byte
}