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
