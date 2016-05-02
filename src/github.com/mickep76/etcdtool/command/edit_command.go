package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/iodatafmt"
	"golang.org/x/net/context"
)

// NewEditCommand sets data from input.
func NewEditCommand() cli.Command {
	return cli.Command{
		Name:  "edit",
		Usage: "edit a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort, s", Usage: "returns result in sorted order"},
			cli.BoolFlag{Name: "yes, y", Usage: "Answer yes to any questions"},
			cli.BoolFlag{Name: "replace, r", Usage: "Replace data"},
			cli.BoolFlag{Name: "validate, v", EnvVar: "ETCDTOOL_VALIDATE", Usage: "Validate data before import"},
			cli.StringFlag{Name: "format, f", Value: "JSON", EnvVar: "ETCDTOOL_FORMAT", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "editor, e", Value: "vim", Usage: "Editor", EnvVar: "EDITOR"},
			cli.StringFlag{Name: "tmp-file, t", Value: ".etcdtool", Usage: "Temporary file"},
		},
		Action: func(c *cli.Context) error {
			editCommandFunc(c)
			return nil
		},
	}
}

func editFile(editor string, file string) error {
	_, err := exec.LookPath(editor)
	if err != nil {
		fatalf("Editor doesn't exist: %s", editor)
	}

	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fatal(err.Error())
	}
	return nil
}

// editCommandFunc edit data as either JSON, YAML or TOML.
func editCommandFunc(c *cli.Context) {
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

	sort := c.Bool("sort")

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		fatal(err.Error())
	}

	// Temporary file append file type to support syntax highlighting
	tmpfile := c.String("tmp-file") + "." + strings.ToLower(c.String("format"))

	// Check if dir exists and is a directory.
	exists, err := dirExists(dir, c, ki)
	if err != nil {
		fatal(err.Error())
	}

	if !exists {
		if askYesNo(fmt.Sprintf("Dir. doesn't exist: %s create it", dir)) {
			// Create dir.
			if _, err := ki.Set(context.TODO(), dir, "", &client.SetOptions{Dir: true}); err != nil {
				fatal(err.Error())
			}
			exists = true
		} else {
			os.Exit(1)
		}
	}

	// If file exist's resume editing?
	if _, err := os.Stat(tmpfile); os.IsNotExist(err) {
		// Export to file.
		exportFunc(dir, sort, tmpfile, f, c, ki)
	} else {
		if !askYesNo(fmt.Sprintf("Temp. file already exist's resume editing")) {
			// Export to file.
			exportFunc(dir, sort, tmpfile, f, c, ki)
		}
	}

	// Get modified time stamp.
	before, err := os.Stat(tmpfile)
	if err != nil {
		fatal(err.Error())
	}

	// Edit file.
	editFile(c.String("editor"), tmpfile)

	// Check modified time stamp.
	after, err := os.Stat(tmpfile)
	if err != nil {
		fatal(err.Error())
	}

	// Import from file if it has changed.
	if before.ModTime() != after.ModTime() {
		importFunc(dir, tmpfile, f, c.Bool("replace"), c.Bool("yes"), e, c, ki)
	} else {
		fmt.Printf("File wasn't modified, skipping import\n")
	}

	// Unlink file.
	if err := os.Remove(tmpfile); err != nil {
		fatal(err.Error())
	}
}
