package command

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
	"golang.org/x/net/context"
)

// NewImportCommand sets data from input.
func NewImportCommand() cli.Command {
	return cli.Command{
		Name:  "import",
		Usage: "import a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "yes, y", Usage: "Answer yes to any questions"},
			cli.BoolFlag{Name: "replace, r", Usage: "Replace data"},
			cli.BoolFlag{Name: "validate, v", EnvVar: "ETCDTOOL_VALIDATE", Usage: "Validate data before import"},
			cli.StringFlag{Name: "format, f", Value: "JSON", EnvVar: "ETCDTOOL_FORMAT", Usage: "Data serialization format YAML, TOML or JSON"},
		},
		Action: func(c *cli.Context) error {
			importCommandFunc(c)
			return nil
		},
	}
}

func dirExists(dir string, c *cli.Context, ki client.KeysAPI) (bool, error) {
	ctx, cancel := contextWithCommandTimeout(c)
	_, err := ki.Get(ctx, dir, &client.GetOptions{})
	cancel()
	if err != nil {
		if cerr, ok := err.(client.Error); ok && cerr.Code == 100 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isDir(dir string, c *cli.Context, ki client.KeysAPI) (bool, error) {
	ctx, cancel := contextWithCommandTimeout(c)
	resp, err := ki.Get(ctx, dir, &client.GetOptions{})
	cancel()
	if err != nil {
		return false, err
	}
	if resp.Node.Dir {
		return false, nil
	}
	return true, nil
}

func askYesNo(msg string) bool {
	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [yes/no]? ", msg)
		inp, _, err := stdin.ReadLine()
		if err != nil {
			fatal(err.Error())
		}

		switch strings.ToLower(string(inp)) {
		case "yes":
			return true
		case "no":
			return false
		default:
			fmt.Printf("Incorrect input: %s\n ", inp)
		}
	}
}

// importCommandFunc imports data as either JSON, YAML or TOML.
func importCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		fatal("You need to specify directory")
	}
	dir := c.Args()[0]

	// Remove trailing slash.
	if dir != "/" {
		dir = strings.TrimRight(dir, "/")
	}
	infof("Using dir: %s", dir)

	if len(c.Args()) == 1 {
		fatal("You need to specify input file")
	}
	input := c.Args()[1]

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		fatal(err.Error())
	}

	// Load configuration file.
	e := loadConfig(c)

	// New dir API.
	ki := newKeyAPI(e)

	importFunc(dir, input, f, c.Bool("replace"), c.Bool("yes"), e, c, ki)
}

func importFunc(dir string, file string, f iodatafmt.DataFmt, replace bool, yes bool, e Etcdtool, c *cli.Context, ki client.KeysAPI) {
	// Check if dir exists and is a directory.
	exists, err := dirExists(dir, c, ki)
	if err != nil {
		fatalf("Specified dir doesn't exist: %s", dir)
	}

	if exists {
		exist, err := isDir(dir, c, ki)
		if err != nil {
			fatal(err.Error())
		}

		if exist {
			fatalf("Specified dir is not a directory: %s", dir)
		}
	}

	// Load file.
	m, err := iodatafmt.Load(file, f)
	if err != nil {
		fatal(err.Error())
	}

	// Validate data.
	if c.Bool("validate") {
		validateFunc(e, dir, m)
	}

	if exists {
		if replace {
			if !yes {
				if !askYesNo(fmt.Sprintf("Do you want to overwrite data in directory: %s", dir)) {
					os.Exit(1)
				}
			}

			// Delete dir.
			if _, err = ki.Delete(context.TODO(), dir, &client.DeleteOptions{Recursive: true}); err != nil {
				fatal(err.Error())
			}
		} else {
			if !yes {
				if !askYesNo(fmt.Sprintf("Do you want to overwrite data in directory: %s", dir)) {
					os.Exit(1)
				}
			}
		}
	} else {
		// Create dir.
		if _, err := ki.Set(context.TODO(), dir, "", &client.SetOptions{Dir: true}); err != nil {
			fatal(err.Error())
		}
	}

	// Import data.
	if err = etcdmap.Create(ki, dir, reflect.ValueOf(m)); err != nil {
		fatal(err.Error())
	}
}
