package utils

import "time"

const (
	// GameTick - 60 updates per second
	GameTick = time.Second / 60

	// RefreshTick - 20 frames per second
	RefreshTick = time.Second / 20

	// Board and graphics
	BoardMaxX = 111
	BoardMaxY =  30
	BoardBuffer  = 5
	PaddleLength = 6
)
