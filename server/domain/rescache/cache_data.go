package rescache

import "time"

type CacheData struct {
	Code         string    `json:"code"`
	User         string    `json:"user"`
	Response     string    `json:"response"`
	ResponseCode string    `json:"response_code"`
	Error        string    `json:"error"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Progress     string    `json:"progress"`
}
