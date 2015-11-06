package command

import (
	"fmt"
	"io/ioutil"
	"log"
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
			cli.BoolFlag{Name: "yes", Usage: "Answer yes to any questions"},
			cli.BoolFlag{Name: "replace", Usage: "Delete entry before import"},
			cli.StringFlag{Name: "format", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "input", Value: "JSON", Usage: "Input File"},
		},
		Action: func(c *cli.Context) {
			exportCommandFunc(c, mustNewKeyAPI(c))
		},
	}
}

// importCommandFunc imports data as either JSON, YAML or TOML.
func importCommandFunc(c *cli.Context, ki client.KeysAPI) {
	var key string
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify directory")
	} else {
		key = c.Args()[0]
	}

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Import data.
	fi, _ := os.Stdin.Stat()
	var m interface{}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		b, _ := ioutil.ReadAll(os.Stdin)
		var err error
		m, err = iodatafmt.Unmarshal(b, f)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else if c.String("input") != "" {
		var err error
		m, err = iodatafmt.Load(c.String("input"), f)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Fatal("No input provided")
	}

	// Delete dir.
	if c.Bool("replace") {
		if !c.Bool("yes") {
			fmt.Printf("Do you want to replace data in directory: %s: %s? [yes|no]", strings.TrimRight(key, "/"))
			var query string
			fmt.Scanln(&query)
			if strings.ToLower(query) != "yes" {
				os.Exit(0)
			}
		}

		if _, err = ki.Delete(context.TODO(), strings.TrimRight(key, "/"), &client.DeleteOptions{Recursive: true}); err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Check if directory exist's

	// Create dir.
	if _, err := ki.Set(context.TODO(), key, "", &client.SetOptions{Dir: true}); err != nil {
		log.Printf(err.Error())

		// Should prob. check that we're actually dealing with an existing key and not something else...
		fmt.Printf("Do you want to overwrite data in directory: %s? [yes|no]", strings.TrimRight(key, "/"))
		var query string
		fmt.Scanln(&query)
		if strings.ToLower(query) != "yes" {
			os.Exit(0)
		}
	}

	// Import data.
	if err = etcdmap.Create(ki, strings.TrimRight(key, "/"), reflect.ValueOf(m)); err != nil {
		log.Fatal(err.Error())
	}
}
