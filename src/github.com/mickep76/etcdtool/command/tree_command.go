package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func NewTreeCommand() cli.Command {
	return cli.Command{
		Name:  "tree",
		Usage: "List directory as a tree",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort", Usage: "returns result in sorted order"},
		},
		Action: func(c *cli.Context) {
			treeCommandFunc(c, newKeyAPI(c))
		},
	}
}

var numDirs int
var numKeys int

// treeCommandFunc executes the "tree" command.
func treeCommandFunc(c *cli.Context, ki client.KeysAPI) {
	key := "/"
	if len(c.Args()) != 0 {
		key = c.Args()[0]
	}

	sort := c.Bool("sort")

	resp, err := ki.Get(context.TODO(), key, &client.GetOptions{Sort: sort, Recursive: true})
	if err != nil {
		log.Fatal(err.Error())
	}

	numDirs = 0
	numKeys = 0
	fmt.Println(strings.TrimRight(key, "/") + "/")
	printTree(resp.Node, "")
	fmt.Printf("\n%d directories, %d keys\n", numDirs, numKeys)
}

// printTree writes a response out in a manner similar to the `tree` command in unix.
func printTree(root *client.Node, indent string) {
	for i, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]

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
