package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	jsonschema "github.com/xeipuuv/gojsonschema"

	"github.com/mickep76/common"
)

func main() {
	// Get connection env variable.
	conn := common.GetEnv()

	// Options.
	version := flag.Bool("version", false, "Version")
	node := flag.String("node", "", "Etcd node")
	port := flag.String("port", "2379", "Etcd port")
	dir := flag.String("dir", "", "Etcd directory")
	schema := flag.String("schema", "", "Etcd key for JSON schema")
	flag.Parse()

	// Print version.
	if *version {
		fmt.Printf("etcd-validate %s\n", common.Version)
		os.Exit(0)
	}

	// Validate input.
	if len(conn) < 1 && *node == "" {
		log.Fatalf("You need to specify Etcd host.")
	}

	if *dir == "" {
		log.Fatalf("You need to specify Etcd dir.")
	}

	if *schema == "" {
		log.Fatalf("You need to specify Etcd key for JSON schema.")
	}

	// Setup Etcd client.
	if *node != "" {
		conn = []string{fmt.Sprintf("http://%v:%v", *node, *port)}
	}
	client := etcd.NewClient(conn)

	// Get JSON Schema.
	res, err := client.Get(*schema, false, false)
	if err != nil {
		log.Fatal(err.Error())
	}
	schemaLoader := jsonschema.NewStringLoader(res.Node.Value)

	// Get Etcd dir.
	res2, err2 := client.Get(*dir, true, true)
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	d := etcdmap.Map(res2.Node)

	for e, v := range d {
		docLoader := jsonschema.NewGoLoader(v)

		result, err := jsonschema.Validate(schemaLoader, docLoader)
		if err != nil {
			panic(err.Error())
		}

		if !result.Valid() {
			fmt.Printf("### %s ###\n", e)
			for _, e := range result.Errors() {
				fmt.Printf("- %s: %s\n", e.Field(), e.Description())
			}
		}
	}
}
