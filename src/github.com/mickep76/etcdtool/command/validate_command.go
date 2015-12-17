package command

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/xeipuuv/gojsonschema"
)

// NewValidateCommand sets data from input.
func NewValidateCommand() cli.Command {
	return cli.Command{
		Name:  "validate",
		Usage: "validate a directory",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			validateCommandFunc(c)
		},
	}
}

// validateCommandFunc validate data using JSON Schema.
func validateCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify directory")
	}
	dir := c.Args()[0]

	// Remove trailing slash.
	if dir != "/" {
		dir = strings.TrimRight(dir, "/")
	}
	Infof(c, "Using dir: %s", dir)

	// Load configuration file.
	e := LoadConfig(c)

	// New dir API.
	ki := newKeyAPI(e)

	// Map directory to routes.
	var schema string
	for _, r := range e.Routes {
		match, err := regexp.MatchString(r.Regexp, dir)
		if err != nil {
			log.Fatal(err.Error())
		}
		if match {
			schema = r.Schema
		}
	}

	if schema == "" && len(c.Args()) == 1 {
		log.Fatal("You need to specify JSON schema URI")
	}

	if len(c.Args()) > 1 {
		schema = c.Args()[1]
	}

	// Get directory.
	ctx, cancel := contextWithCommandTimeout(c)
	resp, err := ki.Get(ctx, dir, &client.GetOptions{Recursive: true})
	cancel()
	if err != nil {
		log.Fatal(err.Error())
	}
	m := etcdmap.Map(resp.Node)

	// Validate directory.
	Infof(c, "Using JSON schema: %s", schema)
	schemaLoader := gojsonschema.NewReferenceLoader(schema)
	docLoader := gojsonschema.NewGoLoader(m)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Print results.
	if !result.Valid() {
		for _, e := range result.Errors() {
			fmt.Printf("%s: %s\n", strings.Replace(e.Context().String("/"), "(root)", dir, 1), e.Description())
		}
	}
}
