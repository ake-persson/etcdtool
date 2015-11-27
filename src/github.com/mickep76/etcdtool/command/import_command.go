package command

import (
	"bufio"
	"fmt"
	"log"
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
			cli.StringFlag{Name: "format, f", Value: "JSON", EnvVar: "ETCDTOOL_FORMAT", Usage: "Data serialization format YAML, TOML or JSON"},
		},
		Action: func(c *cli.Context) {
			importCommandFunc(c, newKeyAPI(c))
		},
	}
}

func keyExists(key string, c *cli.Context, ki client.KeysAPI) (bool, error) {
	ctx, cancel := contextWithCommandTimeout(c)
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
	ctx, cancel := contextWithCommandTimeout(c)
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
		fmt.Printf("%s [yes/no]? ", msg)
		inp, _, err := stdin.ReadLine()
		if err != nil {
			log.Fatal(err.Error())
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
func importCommandFunc(c *cli.Context, ki client.KeysAPI) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify directory")
	}
	// Fix for root
	key := strings.TrimRight(c.Args()[0], "/") //+ "/"

	if len(c.Args()) == 1 {
		log.Fatal("You need to specify input file")
	}
	input := c.Args()[1]

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		log.Fatal(err.Error())
	}

	importFunc(key, input, f, c.Bool("replace"), c.Bool("yes"), c, ki)
}

func importFunc(key string, file string, f iodatafmt.DataFmt, replace bool, yes bool, c *cli.Context, ki client.KeysAPI) {
	// Check if key exists and is a directory.
	exists, err := keyExists(key, c, ki)
	if err != nil {
		log.Fatalf("Specified key doesn't exist: %s", key)
	}

	if exists {
		dir, err := isDir(key, c, ki)
		if err != nil {
			log.Fatal(err.Error())
		}

		if dir {
			log.Fatalf("Specified key is not a directory: %s", key)
		}
	}

	// Load file.
	m, err := iodatafmt.Load(file, f)
	if err != nil {
		log.Fatal(err.Error())
	}

	if exists {
		if replace {
			if !askYesNo(fmt.Sprintf("Do you want to overwrite data in directory: %s", key)) {
				os.Exit(1)
			}

			// Delete dir.
			if _, err = ki.Delete(context.TODO(), key, &client.DeleteOptions{Recursive: true}); err != nil {
				log.Fatal(err.Error())
			}
		} else {
			if !yes {
				if !askYesNo(fmt.Sprintf("Do you want to overwrite data in directory: %s", key)) {
					os.Exit(1)
				}
			}
		}
	} else {
		// Create dir.
		if _, err := ki.Set(context.TODO(), key, "", &client.SetOptions{Dir: true}); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Import data.
	if err = etcdmap.Create(ki, key, reflect.ValueOf(m)); err != nil {
		log.Fatal(err.Error())
	}
}
