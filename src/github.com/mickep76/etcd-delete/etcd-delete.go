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
	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	force := flag.Bool("force", false, "Force delete without asking")
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
