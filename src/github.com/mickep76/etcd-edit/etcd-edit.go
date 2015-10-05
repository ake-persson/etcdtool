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
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
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
	new := flag.Bool("new", false, "Create new directory entry using template")
	template := flag.String("template", "", "etcd key for template")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	force := flag.Bool("force", false, "Force delete without asking")
	noDelete := flag.Bool("no-delete", false, "Don't delete entry before import")
	noValidate := flag.Bool("no-validate", false, "Skip validation using JSON schema")
	schema := flag.String("schema", "", "etcd key for JSON schema")
	editor := flag.String("editor", getEditor(), "Editor")
	dir := flag.String("dir", "/", "etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	tmpFile := flag.String("tmp-file", ".etcd-edit.swp", "Temporary file")
	flag.Parse()

	// Print version.
	if *version {
		fmt.Printf("etcd-edit %s\n", common.Version)
		os.Exit(0)
	}

	// Get data format.
	f, err := iodatafmt.Format(*format)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Connect to etcd.
	cfg := etcd.Config{
		Endpoints:               strings.Split(*peers, ","),
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	client, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !*noValidate && *schema == "" {
		// Get routes.
		kapi := etcd.NewKeysAPI(client)
		res, err := kapi.Get(context.Background(), "/routes", &etcd.GetOptions{Recursive: true})
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
					*template = vm["template"].(string)
					break
				}
			}
		}

		if *schema == "" {
			log.Fatalf("Couldn't determine schema and template to use for directory (use -no-validate to skip this): %s", *dir)
		}
	}

	var m map[string]interface{}

	// Export data.
	if *new {
		// Get JSON template.
		kapi := etcd.NewKeysAPI(client)
		res, err := kapi.Get(context.Background(), *template, &etcd.GetOptions{})
		if err != nil {
			log.Fatal(err.Error())
		}
		m2, err2 := iodatafmt.Unmarshal([]byte(res.Node.Value), iodatafmt.JSON)
		if err2 != nil {
			log.Fatal(err2.Error())
		}
		m = m2.(map[string]interface{})

	} else {
		kapi := etcd.NewKeysAPI(client)
		res, err := kapi.Get(context.Background(), *dir, &etcd.GetOptions{Recursive: true})
		if err != nil {
			log.Fatal(err.Error())
		}
		m = etcdmap.Map(res.Node)
	}

	// Write output.
	iodatafmt.Write(*tmpFile, m, f)

EDIT:

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
		log.Printf(err4.Error())

		fmt.Printf("Do you want to correct the changes? [yes|no]")
		var query string
		fmt.Scanln(&query)
		if query != "yes" {
			os.Exit(0)
		}

		goto EDIT
	}

	// Validate input.
	if !*noValidate {

		// Get JSON Schema.
		kapi := etcd.NewKeysAPI(client)
		res, err := kapi.Get(context.Background(), *schema, &etcd.GetOptions{})
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
				fmt.Printf("%s: %s\n", strings.Replace(e.Context().String("/"), "(root)", *dir, 1), e.Description())
			}

			fmt.Printf("Do you want to correct the changes? [yes|no]")
			var query string
			fmt.Scanln(&query)
			if query != "yes" {
				os.Exit(0)
			}

			goto EDIT
		}
	}

	// TODO: Properly check if key does exist
	// Delete dir.
	if !*noDelete {
		if !*force {
			fmt.Printf("Do you want to remove existing data in path: %s? [yes|no]", strings.TrimRight(*dir, "/"))
			var query string
			fmt.Scanln(&query)
			if query != "yes" {
				os.Exit(0)
			}
		}

		kapi := etcd.NewKeysAPI(client)
		if _, err = kapi.Delete(context.Background(), strings.TrimRight(*dir, "/"), &etcd.DeleteOptions{Recursive: true}); err != nil {
			// Don't exit since -new dir. won't exist
			log.Println(err.Error())
		} else {
			log.Printf("Removed path: %s", strings.TrimRight(*dir, "/"))
		}

		// Create dir.
		if _, err := kapi.Set(context.TODO(), *dir, "", &etcd.SetOptions{Dir: true}); err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		fmt.Printf("Do you want to overwrite existing data in path: %s? [yes|no]", strings.TrimRight(*dir, "/"))
		var query string
		fmt.Scanln(&query)
		if query != "yes" {
			os.Exit(0)
		}
	}

	// Import data.
	if err = etcdmap.Create(&client, strings.TrimRight(*dir, "/"), reflect.ValueOf(m)); err != nil {
		log.Fatal(err.Error())
	}
}
