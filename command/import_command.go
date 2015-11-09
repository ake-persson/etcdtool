package command

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
)

func NewImportCommand() cli.Command {
	return cli.Command{
		Name:  "import",
		Usage: "import a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "yes, y", Usage: "Answer yes to any questions"},
			cli.BoolFlag{Name: "replace, r", Usage: "Replace data"},
			cli.StringFlag{Name: "format, f", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "input, i", Value: "", Usage: "Input File"},
		},
		Action: func(c *cli.Context) {
			importCommandFunc(c, mustNewKeyAPI(c))
		},
	}
}

func keyExists(key string, c *cli.Context, ki client.KeysAPI) (bool, error) {
	ctx, cancel := contextWithTotalTimeout(c)
	_, err := ki.Get(ctx, key, &client.GetOptions{})
	cancel()
	if err != nil {
		if cerr, ok := err.(client.Error); ok && cerr.Code == 100 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isDir(key string, c *cli.Context, ki client.KeysAPI) (bool, error) {
	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Get(ctx, key, &client.GetOptions{})
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
		inp, _, err := stdin.ReadLine()
		if err != nil {
			handleError(ExitServerError, err)
		}

		switch strings.ToLower(string(inp)) {
		case "yes":
			return true
		case "no":
			return false
		default:
			fmt.Printf("Incorrect input: %s\n", inp)
		}
	}
}

// importCommandFunc imports data as either JSON, YAML or TOML.
func importCommandFunc(c *cli.Context, ki client.KeysAPI) {
	var key string
	if len(c.Args()) == 0 {
		handleError(ExitServerError, errors.New("You need to specify directory"))
	} else {
		key = c.Args()[0]
	}

	// Check if key exists and is a directory.
	exists, err := keyExists(key, c, ki)
	if err != nil {
		handleError(ExitServerError, errors.New("No input provided"))
	}

	if exists {
		dir, err := isDir(key, c, ki)
		if err != nil {
			handleError(ExitServerError, err)
		}

		if dir {
			handleError(ExitServerError, fmt.Errorf("Specified key is not a directory: %s", key))
		}
	}

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		handleError(ExitServerError, errors.New("No input provided"))
	}

	// Import data.
	if c.String("input") == "" {
		handleError(ExitServerError, errors.New("No input provided"))
	}

	m, err := iodatafmt.Load(c.String("input"), f)
	if err != nil {
		handleError(ExitServerError, err)
	}

	if exists {
		if c.Bool("replace") {
			if !askYesNo(fmt.Sprintf("Do you want to overwrite data in directory: %s", strings.TrimRight(key, "/"))) {
				os.Exit(1)
			}

			// Delete dir.
			if _, err = ki.Delete(context.TODO(), strings.TrimRight(key, "/"), &client.DeleteOptions{Recursive: true}); err != nil {
				handleError(ExitServerError, err)
			}
		} else {
			if !c.Bool("yes") {
				if !askYesNo(fmt.Sprintf("Do you want to overwrite data in directory: %s", strings.TrimRight(key, "/"))) {
					os.Exit(1)
				}
			}
		}
	} else {
		// Create dir.
		if _, err := ki.Set(context.TODO(), key, "", &client.SetOptions{Dir: true}); err != nil {
			handleError(ExitServerError, err)
		}
	}

	// Import data.
	if err = etcdmap.Create(ki, strings.TrimRight(key, "/"), reflect.ValueOf(m)); err != nil {
		handleError(ExitServerError, err)
	}
}
