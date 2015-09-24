package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
	jsonschema "github.com/xeipuuv/gojsonschema"

	"github.com/mickep76/common"
)

func getEditor() string {
	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		if a[0] == "EDITOR" {
			return a[1]
		}
	}

	return "vim"
}

func main() {
	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	force := flag.Bool("force", false, "Force delete without asking")
	delete := flag.Bool("delete", false, "Delete entry before import")
	noValidate := flag.Bool("no-validate", false, "Skip validation using JSON schema")
	schema := flag.String("schema", "", "etcd key for JSON schema")
	editor := flag.String("editor", getEditor(), "Editor")
	dir := flag.String("dir", "/", "etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	tmpFile := flag.String("tmp-file", ".etcd-edit.swp", "Temporary file")
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

	// Export data.
	res, err := client.Get(*dir, true, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	m := etcdmap.Map(res.Node)

	// Write output.
	iodatafmt.Write(*tmpFile, m, f)

	_, err2 := exec.LookPath(*editor)
	if err2 != nil {
		log.Fatalf("Editor doesn't exist: %s", *editor)
	}

	cmd := exec.Command(*editor, *tmpFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err3 := cmd.Run()
	if err3 != nil {
		log.Fatalf(err3.Error())
	}

	// Import data.
	var m2 interface{}
	var err4 error
	m2, err4 = iodatafmt.Load(*tmpFile, f)
	if err4 != nil {
		log.Fatal(err4.Error())
	}

	// Validate input.
	if !*noValidate {

		// Get JSON Schema.
		res, err := client.Get(*schema, false, false)
		if err != nil {
			log.Fatal(err.Error())
		}
		schemaLoader := jsonschema.NewStringLoader(res.Node.Value)
		docLoader := jsonschema.NewGoLoader(m2)

		result, err := jsonschema.Validate(schemaLoader, docLoader)
		if err != nil {
			panic(err.Error())
		}

		if !result.Valid() {
			for _, e := range result.Errors() {
				fmt.Printf("  - %s: %s\n", e.Field(), e.Description())
			}
			os.Exit(1)
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

		// Create dir.
		if _, err := client.CreateDir(*dir, 0); err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Import data.
	if err = etcdmap.Create(client, strings.TrimRight(*dir, "/"), reflect.ValueOf(m2)); err != nil {
		log.Fatal(err.Error())
	}
}
