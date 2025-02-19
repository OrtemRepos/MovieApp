package model

type RecordID string
type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserID string
type RatingValue float64

type Rating struct {
	RecordID   RecordID    `json:"record_id"`
	RecordType RecordType  `json:"record_type"`
	UserID     UserID      `json:"user_id"`
	Value      RatingValue `json:"value"`
}

type RatingEvent struct {
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}

type RatingEventType string

const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
