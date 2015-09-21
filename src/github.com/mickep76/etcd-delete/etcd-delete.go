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

func main() {
	// Get connection env variable.
	conn := common.GetEnv()

	// Options.
	version := flag.Bool("version", false, "Version")
	force := flag.Bool("force", false, "Force delete without asking")
	node := flag.String("node", "", "Etcd node")
	port := flag.String("port", "2379", "Etcd port")
	dir := flag.String("dir", "/", "Etcd directory")
	flag.Parse()

	// Print version.
	if *version {
		fmt.Printf("etcd-import %s\n", common.Version)
		os.Exit(0)
	}

	// Validate input.
	if len(conn) < 1 && *node == "" {
		log.Fatalf("You need to specify Etcd host.")
	}

	// Setup Etcd client.
	if *node != "" {
		conn = []string{fmt.Sprintf("http://%v:%v", *node, *port)}
	}
	client := etcd.NewClient(conn)

	if !*force {
		fmt.Printf("Remove path: %s? [yes|no]", strings.TrimRight(*dir, "/"))
		var query string
		fmt.Scanln(&query)
		if query != "yes" {
			os.Exit(0)
		}
	}

	if _, err := client.Delete(strings.TrimRight(*dir, "/"), true); err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Removed path: %s", strings.TrimRight(*dir, "/"))
}
