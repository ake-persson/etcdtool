package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
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
			validateCommandFunc(c, mustNewKeyAPI(c))
		},
	}
}

// validateCommandFunc validate data using JSON Schema.
func validateCommandFunc(c *cli.Context, ki client.KeysAPI) {
	var key string
	if len(c.Args()) == 0 {
		handleError(ExitServerError, errors.New("You need to specify directory"))
	} else {
		key = strings.TrimRight(c.Args()[0], "/") + "/"
	}

	if len(c.Args()) == 1 {
		handleError(ExitBadArgs, errors.New("You need to specify JSON schema URI"))
	}
	schema := c.Args()[1]

	// Get directory.
	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Get(ctx, key, &client.GetOptions{Recursive: true})
	cancel()
	if err != nil {
		handleError(ExitServerError, err)
	}
	m := etcdmap.Map(resp.Node)

	// Validate directory.
	schemaLoader := gojsonschema.NewReferenceLoader(schema)
	docLoader := gojsonschema.NewGoLoader(m)
	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		handleError(ExitServerError, err)
	}

	// Print results.
	if !result.Valid() {
		for _, e := range result.Errors() {
			fmt.Printf("%s: %s\n", strings.Replace(e.Context().String("/"), "(root)", key, 1), e.Description())
		}
	}
}
