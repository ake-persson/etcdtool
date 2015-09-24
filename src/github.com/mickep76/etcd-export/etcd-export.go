package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"

	"github.com/mickep76/common"
)

func main() {
	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	dir := flag.String("dir", "/", "etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	output := flag.String("output", "", "Output file")
	flag.Parse()

	// Print version.
	if *version {
		fmt.Printf("etcd-export %s\n", common.Version)
		os.Exit(0)
	}

	// Get data format.
	f, err := iodatafmt.Format(*format)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Setup etcd client.
	client := etcd.NewClient(strings.Split(*peers, ","))

	// Export data.
	res, err := client.Get(*dir, true, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	m := etcdmap.Map(res.Node)

	// Write output.
	if *output != "" {
		iodatafmt.Write(*output, m, f)
	} else {
		iodatafmt.Print(m, f)
	}
}
