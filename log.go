package main

import (
	"log"
	"os"
)

// Defines custom logger levels.
var (
	Info  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
)

func init() {
	Info = log.New(
		os.Stderr,
		"NFO: ",
		0,
	)

	Error = log.New(
		os.Stderr,
		"ERR: ",
		0,
	)

	Fatal = log.New(
		os.Stderr,
		"FTL: ",
		0,
	)
}
