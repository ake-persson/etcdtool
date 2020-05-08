package command

import (
	"log"
)

var debug bool

func info(msg string) {
	if debug {
		log.Print(msg)
	}
}

func infof(fmt string, args ...interface{}) {
	if debug {
		log.Printf(fmt, args...)
	}
}

func fatal(msg string) {
	log.Fatal(msg)
}

func fatalf(fmt string, args ...interface{}) {
	log.Fatalf(fmt, args...)
}
