package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"

	"github.com/mickep76/common"
)

func main() {
	// Get the FileInfo struct describing the standard input.
	fi, _ := os.Stdin.Stat()

	// Get connection env variable.
	conn := common.GetEnv()

	// Options.
	version := flag.Bool("version", false, "Version")
	node := flag.String("node", "", "Etcd node")
	port := flag.String("port", "2379", "Etcd port")
	dir := flag.String("dir", "/", "Etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	input := flag.String("input", "", "Input file")
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

	// Get data format.
	f, err := iodatafmt.Format(*format)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Setup Etcd client.
	if *node == "" {
		conn = []string{fmt.Sprintf("http://%v:%v", *node, *port)}
	}
	client := etcd.NewClient(conn)

	// Import data.
	m := make(map[string]interface{})
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		b, _ := ioutil.ReadAll(os.Stdin)
		var err error
		m, err = iodatafmt.Unmarshal(b, f)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else if *input != "" {
		var err error
		m, err = iodatafmt.Load(*input, f)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Fatal("No input provided")
	}

	if err = etcdmap.Create(client, strings.TrimRight(*dir, "/"), reflect.ValueOf(m)); err != nil {
		log.Fatal(err.Error())
	}
}
