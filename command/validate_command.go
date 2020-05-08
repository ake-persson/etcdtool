package command

import (
	"fmt"
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
		Action: func(c *cli.Context) error {
			validateCommandFunc(c)
			return nil
		},
	}
}

// validateCommandFunc validate data using JSON Schema.
func validateCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		fatal("You need to specify directory")
	}
	dir := c.Args()[0]

	// Remove trailing slash.
	if dir != "/" {
		dir = strings.TrimRight(dir, "/")
	}
	infof("Using dir: %s", dir)

	// Load configuration file.
	e := loadConfig(c)

	// New dir API.
	ki := newKeyAPI(e)

	// Map directory to routes.
	var schema string
	for _, r := range e.Routes {
		match, err := regexp.MatchString(r.Regexp, dir)
		if err != nil {
			fatal(err.Error())
		}
		if match {
			schema = r.Schema
		}
	}

	if schema == "" && len(c.Args()) == 1 {
		fatal("You need to specify JSON schema URI")
	}

	if len(c.Args()) > 1 {
		schema = c.Args()[1]
	}

	// Get directory.
	ctx, cancel := contextWithCommandTimeout(c)
	resp, err := ki.Get(ctx, dir, &client.GetOptions{Recursive: true})
	cancel()
	if err != nil {
		fatal(err.Error())
	}
	m := etcdmap.Map(resp.Node)

	// Validate directory.
	infof("Using JSON schema: %s", schema)
	schemaLoader := gojsonschema.NewReferenceLoader(schema)
	docLoader := gojsonschema.NewGoLoader(m)
	fmt.Println("Validating...")
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		fatal(fmt.Sprintf("Error attempting to validate: %v", err.Error()))
	}

	// Print results.
	if !result.Valid() {
		for _, err := range result.Errors() {
			fmt.Printf("%s: %s [value: %s]\n", strings.Replace(err.Context().String("/"), "(root)", dir, 1), err.Description(), err.Value())
		}
	}
}

func validateFunc(e Etcdtool, dir string, d interface{}) {
	// Map directory to routes.
	var schema string
	for _, r := range e.Routes {
		match, err := regexp.MatchString(r.Regexp, dir)
		if err != nil {
			fatal(err.Error())
		}
		if match {
			schema = r.Schema
		}
	}

	if schema == "" {
		fatal("Couldn't determine which JSON schema to use for validation")
	}

	/*
	   if schema == "" && len(c.Args()) == 1 {
	       fatal("You need to specify JSON schema URI")
	   }

	   if len(c.Args()) > 1 {
	       schema = c.Args()[1]
	   }
	*/

	// Validate directory.
	infof("Using JSON schema: %s", schema)
	schemaLoader := gojsonschema.NewReferenceLoader(schema)
	docLoader := gojsonschema.NewGoLoader(d)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		fatal(err.Error())
	}

	// Print results.
	if !result.Valid() {
		for _, err := range result.Errors() {
			fmt.Printf("%s: %s [value: %s]\n", strings.Replace(err.Context().String("/"), "(root)", dir, 1), err.Description(), err.Value())
		}
		fatal("Data validation failed aborting")
	}
}
