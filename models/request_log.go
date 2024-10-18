package models

import "time"

type RequestLog struct {
	Timestamp   time.Time         `json:"timestamp"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	FullPath    string            `json:"full_path"`
	GET         map[string]any    `json:"GET"`
	RequestBody map[string]any    `json:"request_body"`
	User        string            `json:"user"`
	StatusCode  int               `json:"status_code"`
	Duration    float64           `json:"duration"`
	Headers     map[string]string `json:"headers"`
}
