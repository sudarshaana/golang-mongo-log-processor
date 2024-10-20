package models

import "time"

type RequestLog struct {
	Timestamp   time.Time      `json:"timestamp" bson:"timestamp"`
	Method      string         `json:"method" bson:"method"`
	Path        string         `json:"path" bson:"path"`
	FullPath    string         `json:"full_path" bson:"full_path"`
	GET         map[string]any `json:"GET" bson:"get"`
	RequestBody map[string]any `json:"request_body" bson:"request_body"`
	User        string         `json:"user" bson:"user"`
	StatusCode  int            `json:"status_code" bson:"status_code"`
	Duration    float64        `json:"duration" bson:"duration"`
	CreatedAt   time.Time
}
