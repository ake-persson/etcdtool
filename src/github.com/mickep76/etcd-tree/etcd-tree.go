package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"

	"github.com/mickep76/common"
)

var numDirs int
var numKeys int

// Tree
func Tree(root *etcd.Node, indent string) {
	for i, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]

		if n.Dir {
			if i == root.Nodes.Len()-1 {
				fmt.Printf("%s└── %s/\n", indent, k)
				Tree(n, indent+"    ")
			} else {
				fmt.Printf("%s├── %s/\n", indent, k)
				Tree(n, indent+"│   ")
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

func main() {
	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	flag.Parse()

	var dir string
	if len(flag.Args()) < 1 {
		dir = "/"
	} else {
		dir = flag.Args()[0]
	}

	// Print version.
	if *version {
		fmt.Printf("etcd-import %s\n", common.Version)
		os.Exit(0)
	}

	// Connect to etcd.
	cfg := etcd.Config{
		Endpoints:               strings.Split(*peers, ","),
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	client, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Export data.
	kapi := etcd.NewKeysAPI(client)
	res, err := kapi.Get(context.Background(), dir, &etcd.GetOptions{Recursive: true, Sort: true})
	if err != nil {
		log.Fatal(err.Error())
	}

	numDirs = 0
	numKeys = 0
	fmt.Println(strings.TrimRight(dir, "/") + "/")
	Tree(res.Node, "")
	fmt.Printf("\n%d directories, %d keys\n", numDirs, numKeys)
}
