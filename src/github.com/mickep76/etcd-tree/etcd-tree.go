package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"

	"github.com/mickep76/common"
)

// Tree
func Tree(root *etcd.Node, indent string) {
	for _, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]
		if n.Dir {
			fmt.Printf("%s├── %s/\n", indent, k)
			Tree(n, indent+"│   ")
		} else {
			fmt.Printf("%s├── %s\n", indent, k)
		}
	}
}

/*
├── postfix
│   ├── LICENSE
│   ├── TLS_LICENSE
│   ├── access
│   ├── aliases
*/

func main() {
	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	dir := flag.String("dir", "", "etcd directory")
	flag.Parse()

	// Print version.
	if *version {
		fmt.Printf("etcd-import %s\n", common.Version)
		os.Exit(0)
	}

	// Validate input.
	if *dir == "" {
		log.Fatalf("You need to specify etcd dir.")
	}

	// Setup etcd client.
	client := etcd.NewClient(strings.Split(*peers, ","))

	// Export data.
	res, err := client.Get(*dir, true, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(strings.TrimRight(*dir, "/") + "/")
	Tree(res.Node, "")
}
