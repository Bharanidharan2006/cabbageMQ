// TODO
// [ ] - Add guards like you can't send payload on QUEUE_EMPTY or HEARTBEAT error message

// --------------------------------------------------
// IMPORTANT
// --------------------------------------------------

// Settle to int32 for everything as it can run on different machines and the protocol may getten broken

// FIXED
// It runs on the assumption, that the machine that I will run this system on is 64 bits, other wise it breaks. Cause I cast
// uint32 to int(should be 64 bits) to accommodate it correctly
// Also this is unlikely to happen as the payload is never going to get that big

package protocol

import (
	"encoding/binary"
	"fmt"
	"math"
)

type MessageType uint8

const (
	PUBLISH MessageType = iota
	PUBLISHED
	QUEUE_FULL
	TOPIC_NOT_FOUND
	REQUEST
	JOB
	JOB_PROCESSING_SUCCESS
	JOB_PROCESSING_FAIL
	QUEUE_EMPTY
	HEARTBEAT
)

const HEADER_SIZE int = 6;

// uint32 works based on the assumption that no payload in message broker should be anywhere near 2gb. If it is then the reference should be passed and not
// the actual data. Add guard rail as well


// Keeping the id as uint64 as it never really exhausts
type Message struct {
	messageType   MessageType
	topic         string
	id            uint64
	payload       []byte
}

// Here we use only []byte and not a pointer as a slice is just a view into an array and it does not hold any data.
// So even when we pass by value, we are still pointing to the same underlying array
func (m *Message) constuctHeader(message []byte,topicLength uint8, payloadLength uint32) {
	message[0] = byte(m.messageType)
	message[1] = byte(topicLength)
	binary.BigEndian.PutUint32(message[2:6], payloadLength) // not a problem as it is always larger than int
}  

func (m *Message) Serialize() ([]byte, error) {
	switch m.messageType {
	case HEARTBEAT, QUEUE_EMPTY, QUEUE_FULL, TOPIC_NOT_FOUND, PUBLISHED:
		message := make([]byte, HEADER_SIZE)
		m.constuctHeader(message, uint8(0), uint32(0))
		return message, nil
	case PUBLISH, JOB:
			if len(m.payload) > math.MaxInt32 {
			return nil, fmt.Errorf("payload: payload size is too large, size: %d", len(m.payload))
		}
		if len(m.topic) > math.MaxUint8 {
			return nil, fmt.Errorf("topic: topic name is too long, length: %d", len(m.topic))
		}
		// Message Types related to job acknowlegments and errors. SHould have the payloadlength as eight
		topicLength := len(m.topic)
		payloadLength := len(m.payload)
		messageSize := HEADER_SIZE + topicLength + payloadLength + 8  // Isnt this gonna cause an overflow?
		message := make([]byte, messageSize)
		m.constuctHeader(message, uint8(topicLength), uint32(payloadLength))
		topicEnd := HEADER_SIZE + topicLength
		copy(message[HEADER_SIZE:topicEnd], []byte(m.topic))
		idEnd := topicEnd + 8
		binary.BigEndian.PutUint64(message[topicEnd:idEnd], m.id)
		payloadEnd := idEnd + len(m.payload) 
		copy(message[idEnd: payloadEnd], m.payload)
		return message, nil
	case JOB_PROCESSING_SUCCESS, JOB_PROCESSING_FAIL:
		if len(m.topic) > math.MaxUint8 {
			return nil, fmt.Errorf("topic: topic name is too long, length: %d", len(m.topic))
		}
		topicLength := len(m.topic)
		messageSize := HEADER_SIZE + topicLength + 8 // Isnt this gonna cause an overflow?
		message := make([]byte, messageSize)
		m.constuctHeader(message, uint8(topicLength), uint32(0))
		topicEnd := HEADER_SIZE + topicLength
		copy(message[HEADER_SIZE:topicEnd], []byte(m.topic))
		idEnd := topicEnd + 8
		binary.BigEndian.PutUint64(message[topicEnd:idEnd], m.id)
		return message, nil
	case REQUEST:
		if len(m.topic) > math.MaxUint8 {
			return nil, fmt.Errorf("topic: topic name is too long, length: %d", len(m.topic))
		}
		topicLength := len(m.topic)
		messageSize := HEADER_SIZE + topicLength // Isnt this gonna cause an overflow?
		message := make([]byte, messageSize)
		m.constuctHeader(message, uint8(topicLength), uint32(0))
		topicEnd := HEADER_SIZE + topicLength
		copy(message[HEADER_SIZE:topicEnd], []byte(m.topic))
		return message, nil
	default:
		return nil, fmt.Errorf("message type: not a valid message type %d", m.messageType)
	}
	
}