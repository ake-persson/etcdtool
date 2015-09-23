package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
	jsonschema "github.com/xeipuuv/gojsonschema"

	"github.com/mickep76/common"
)

func main() {
	// Get the FileInfo struct describing the standard input.
	fi, _ := os.Stdin.Stat()

	// Get connection env variable.
	conn := common.GetEnv()

	// Options.
	version := flag.Bool("version", false, "Version")
	force := flag.Bool("force", false, "Force delete without asking")
	delete := flag.Bool("delete", false, "Delete entry before import")
	node := flag.String("node", "", "Etcd node")
	port := flag.String("port", "2379", "Etcd port")
	dir := flag.String("dir", "", "Etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	input := flag.String("input", "", "Input file")
	noValidate := flag.Bool("no-validate", false, "No validate")
	schema := flag.String("schema", "", "Etcd key for JSON schema (default \"/schemas/<dir>/schema\")")
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

	if *dir == "" {
		log.Fatalf("You need to specify Etcd dir.")
	}

	/*
		if *schema == "" {
			s := fmt.Sprintf("/schemas/%s/schema", *dir)
			schema = &s
		}
	*/

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

	// TODO: Use struct for Routes
	if !*noValidate && *schema == "" {
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
			log.Fatalf("Couldn't determine schema to use for directory (use -no-validate to skip this): %s", *dir)
		}
	}

	// Import data.
	var m interface{}
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

	// Validate input.
	if !*noValidate {

		// Get JSON Schema.
		res, err := client.Get(*schema, false, false)
		if err != nil {
			log.Fatal(err.Error())
		}
		schemaLoader := jsonschema.NewStringLoader(res.Node.Value)
		docLoader := jsonschema.NewGoLoader(m)

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

	// Delete dir.
	if *delete {
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

	if err = etcdmap.Create(client, strings.TrimRight(*dir, "/"), reflect.ValueOf(m)); err != nil {
		log.Fatal(err.Error())
	}
}
