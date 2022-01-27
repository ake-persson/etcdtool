package command

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
)

// NewLsCommand returns keys from directory
func NewLsCommand() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "list a directory",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "output, o", Value: "", Usage: "Output file"},
		},
		Action: func(c *cli.Context) error {
			lsCommandFunc(c)
			return nil
		},
	}
}

// lsCommandFunc does what `etcdctl ls -p` do
func lsCommandFunc(c *cli.Context) {
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

	lsFunc(dir, c.String("output"), c, ki)
}

func lsFunc(dir string, file string, c *cli.Context, ki client.KeysAPI) {
	ctx, cancel := contextWithCommandTimeout(c)
	resp, err := ki.Get(ctx, dir, &client.GetOptions{Sort: true, Recursive: false})
	cancel()
	if err != nil {
		fatal(err.Error())
	}

	m := etcdmap.Map(resp.Node)

	keys := make([]string, 0, len(m))
	for k := range m {
		isDir := false
		switch m[k].(type) {
			case map[string]interface {}:
				isDir = true
			default:
				// nothing to do
		}

		if isDir {
			keys = append(keys, k + "/")
		} else {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var result strings.Builder
	for i, k := range(keys) {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString("/")
		result.WriteString(k)
		i++
	}

	if file != "" {
		file, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fatal(err.Error())
		}
		file.WriteString(result.String())
		file.WriteString("\n")
	} else {
		fmt.Printf("%s\n", result.String())
	}
}

