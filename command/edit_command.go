package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/iodatafmt"
)

// NewImportCommand sets data from input.
func NewEditCommand() cli.Command {
	return cli.Command{
		Name:  "edit",
		Usage: "edit a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort, s", Usage: "returns result in sorted order"},
			cli.BoolFlag{Name: "yes, y", Usage: "Answer yes to any questions"},
			cli.BoolFlag{Name: "replace, r", Usage: "Replace data"},
			cli.StringFlag{Name: "format, f", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "editor, e", Value: "vim", Usage: "Editor"},
			cli.StringFlag{Name: "tmp-file, t", Value: ".etcdfmt.swp", Usage: "Temporary file"},
		},
		Action: func(c *cli.Context) {
			editCommandFunc(c, mustNewKeyAPI(c))
		},
	}
}

func editFile(editor string, file string) error {
	_, err := exec.LookPath(editor)
	if err != nil {
		handleError(ExitServerError, fmt.Errorf("Editor doesn't exist: %s", editor))
	}

	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		handleError(ExitServerError, err)
	}
	return nil
}

// editCommandFunc edit data as either JSON, YAML or TOML.
func editCommandFunc(c *cli.Context, ki client.KeysAPI) {
	var key string
	if len(c.Args()) == 0 {
		handleError(ExitServerError, errors.New("You need to specify directory"))
	} else {
		key = strings.TrimRight(c.Args()[0], "/") + "/"
	}

	sort := c.Bool("sort")

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		handleError(ExitServerError, err)
	}

	// Export to file.
	exportFunc(key, sort, c.String("tmp-file"), f, c, ki)

	// Edit file.
	editFile(c.String("editor"), c.String("tmp-file"))

	// Check if file changed modified time...

	// Import from file.
	importFunc(key, c.String("tmp-file"), f, c.Bool("replace"), c.Bool("yes"), c, ki)

	// Unlink file.
	if err := os.Remove(c.String("tmp-file")); err != nil {
		handleError(ExitServerError, err)
	}
}
