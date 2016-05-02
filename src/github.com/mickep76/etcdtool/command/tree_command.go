package command

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// NewTreeCommand print directory as a tree.
func NewTreeCommand() cli.Command {
	return cli.Command{
		Name:  "tree",
		Usage: "List directory as a tree",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort", Usage: "returns result in sorted order"},
		},
		Action: func(c *cli.Context) error {
			treeCommandFunc(c)
			return nil
		},
	}
}

var numDirs int
var numKeys int

// treeCommandFunc executes the "tree" command.
func treeCommandFunc(c *cli.Context) {
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

	resp, err := ki.Get(context.TODO(), dir, &client.GetOptions{Sort: sort, Recursive: true})
	if err != nil {
		fatal(err.Error())
	}

	numDirs = 0
	numKeys = 0
	fmt.Println(strings.TrimRight(dir, "/") + "/")
	printTree(resp.Node, "")
	fmt.Printf("\n%d directories, %d dirs\n", numDirs, numKeys)
}

// printTree writes a response out in a manner similar to the `tree` command in unix.
func printTree(root *client.Node, indent string) {
	for i, n := range root.Nodes {
		dirs := strings.Split(n.Key, "/")
		k := dirs[len(dirs)-1]

		if n.Dir {
			if i == root.Nodes.Len()-1 {
				fmt.Printf("%s└── %s/\n", indent, k)
				printTree(n, indent+"    ")
			} else {
				fmt.Printf("%s├── %s/\n", indent, k)
				printTree(n, indent+"│   ")
			}
			numDirs++
		} else {
			if i == root.Nodes.Len()-1 {
				fmt.Printf("%s└── %s\n", indent, k)
			} else {
				fmt.Printf("%s├── %s\n", indent, k)
			}

			numKeys++
		}
	}
}
