package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	jsonschema "github.com/xeipuuv/gojsonschema"

	"github.com/mickep76/common"
)

type Route struct {
	Path   string `json:"path"`
	Regexp string `json:"regexp"`
	Desc   string `json:"desc"`
	Schema string `json:"schema"`
}

func main() {
	// Get connection env variable.
	conn := common.GetEnv()

	// Options.
	version := flag.Bool("version", false, "Version")
	node := flag.String("node", "", "Etcd node")
	port := flag.String("port", "2379", "Etcd port")
	dir := flag.String("dir", "", "Etcd directory")
	schema := flag.String("schema", "", "Etcd key for JSON schema (default \"/schemas/<dir>/schema\")")
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

	/*
		if *schema == "" {
			s := fmt.Sprintf("/schemas/%s/schema", *dir)
			schema = &s
		}
	*/

	// TODO: Check dir is a dir and not a key

	// Setup Etcd client.
	if *node != "" {
		conn = []string{fmt.Sprintf("http://%v:%v", *node, *port)}
	}
	client := etcd.NewClient(conn)

	// TODO: Use struct for Routes
	if *schema == "" {
		// Get routes.
		res, err := client.Get("/routes", true, true)
		if err != nil {
			log.Fatal(err.Error())
		}
		routes := etcdmap.Map(res.Node)

		for _, v := range routes {
			switch reflect.ValueOf(v).Kind() {
			case reflect.Map:
				var vm map[string]interface{}
				vm = v.(map[string]interface{})
				match, _ := regexp.MatchString(vm["regexp"].(string), *dir)
				if match {
					*schema = vm["schema"].(string)
					break
				}
			}
		}

		if *schema == "" {
			log.Fatalf("Couldn't determine schema to use for directory: %s", *dir)
		}

		fmt.Println(*schema)
	}

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
			fmt.Printf("\n%s/%s\n", *dir, e)
			for _, e := range result.Errors() {
				fmt.Printf("  - %s: %s\n", e.Field(), e.Description())
			}
		}
	}
}
