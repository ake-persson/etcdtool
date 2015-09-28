package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

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
	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	dir := flag.String("dir", "", "etcd directory")
	schema := flag.String("schema", "", "etcd key for JSON schema")
	flag.Parse()

	// Print version.
	if *version {
		fmt.Printf("etcd-validate %s\n", common.Version)
		os.Exit(0)
	}

	// Validate input.
	if *dir == "" {
		log.Fatalf("You need to specify etcd dir.")
	}

	// Setup etcd client.
	client := etcd.NewClient(strings.Split(*peers, ","))

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
				match, err := regexp.MatchString(vm["regexp"].(string), *dir)
				if err != nil {
					panic(err)
				}
				if match {
					*schema = vm["schema"].(string)
					break
				}
			}
		}

		if *schema == "" {
			log.Fatalf("Couldn't determine schema to use for directory: %s", *dir)
		}
	}

	// Get JSON Schema.
	res, err := client.Get(*schema, false, false)
	if err != nil {
		log.Fatal(err.Error())
	}
	schemaLoader := jsonschema.NewStringLoader(res.Node.Value)

	// Get etcd dir.
	res2, err2 := client.Get(*dir, true, true)
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	d := etcdmap.Map(res2.Node)

	docLoader := jsonschema.NewGoLoader(d)

	result, err := jsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		panic(err.Error())
	}

	if !result.Valid() {
		for _, e := range result.Errors() {
			fmt.Printf("  - %s: %s\n", e.Field(), e.Description())
		}
	}
}
