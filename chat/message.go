package main

import "time"

type message struct {
	Email   string
	Message string
	When    time.Time
}
