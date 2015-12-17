package command

import (
	"log"

	"github.com/codegangsta/cli"
)

func info(c *cli.Context, msg string) {
	if c.GlobalBool("debug") {
		log.Print(msg)
	}
}

func infof(c *cli.Context, fmt string, args ...interface{}) {
	if c.GlobalBool("debug") {
		log.Printf(fmt, args...)
	}
}

func fatal(c *cli.Context, msg string) {
	log.Fatal(msg)
}

func fatalf(c *cli.Context, fmt string, args ...interface{}) {
	log.Fatalf(fmt, args...)
}
