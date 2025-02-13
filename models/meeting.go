package models

import "gorm.io/gorm"

// ZoomMeeting represents the meeting structure
type ZoomMeeting struct {
	gorm.Model
	ZoomID    int    `json:"zoom_id"`
	Topic     string `json:"topic"`
	JoinURL   string `json:"join_url"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
	Password  string `json:"password,omitempty"`
}
