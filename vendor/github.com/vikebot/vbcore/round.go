package vbcore

import "time"

const (
	RoundStatusOpen     = 1
	RoundStatusClosed   = 2
	RoundStatusRunning  = 3
	RoundStatusFinished = 4
)

type Round struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Wallpaper   string    `json:"wallpaper"`
	Joined      int       `json:"joined"`
	Min         int       `json:"min"`
	Max         int       `json:"max"`
	Starttime   time.Time `json:"starttime"`
	RoundStatus int       `json:"status"`
}
