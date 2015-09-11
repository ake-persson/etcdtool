package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"

	"github.com/mickep76/common"
)

func main() {
	// Get connection env variable.
	conn := common.GetEnv()

	// Options.
	version := flag.Bool("version", false, "Version")
	node := flag.String("node", "", "Etcd node")
	port := flag.String("port", "2379", "Etcd port")
	dir := flag.String("dir", "/", "Etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	output := flag.String("output", "", "Output file")
	flag.Parse()

	// Print version.
	if *version {
		fmt.Printf("etcd-export %s\n", common.Version)
		os.Exit(0)
	}

	// Validate input.
	if len(conn) < 1 && *node == "" {
		log.Fatalf("You need to specify Etcd host.")
	}

	// Get data format.
	f, err := iodatafmt.Format(*format)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Setup Etcd client.
	if *node != "" {
		conn = []string{fmt.Sprintf("http://%v:%v", *node, *port)}
	}
	client := etcd.NewClient(conn)

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
