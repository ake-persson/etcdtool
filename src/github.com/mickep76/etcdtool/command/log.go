package command

import (
	"log"

	"github.com/codegangsta/cli"
)

// Info log.
func Info(c *cli.Context, msg string) {
	if c.GlobalBool("debug") {
		log.Print(msg)
	}
}

// Infof log.
func Infof(c *cli.Context, fmt string, args ...interface{}) {
	if c.GlobalBool("debug") {
		log.Printf(fmt, args...)
	}
}

// Fatal log.
func Fatal(c *cli.Context, msg string) {
	log.Fatal(msg)
}

// Fatalf log.
func Fatalf(c *cli.Context, fmt string, args ...interface{}) {
	log.Fatalf(fmt, args...)
}
