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
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
	//	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
	jsonschema "github.com/xeipuuv/gojsonschema"

	"github.com/mickep76/common"
)

func main() {
	// Get the FileInfo struct describing the standard input.
	fi, _ := os.Stdin.Stat()

	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	force := flag.Bool("force", false, "Force delete without asking")
	delete := flag.Bool("delete", false, "Delete entry before import")
	dir := flag.String("dir", "", "etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	input := flag.String("input", "", "Input file")
	noValidate := flag.Bool("no-validate", false, "Skip validation using JSON schema")
	schema := flag.String("schema", "", "etcd key for JSON schema")
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
		kapi := etcd.NewKeysAPI(client)
		res, err := kapi.Get(context.Background(), *schema, &etcd.GetOptions{})
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
				fmt.Printf("%s: %s\n", strings.Replace(e.Context().String("/"), "(root)", *dir, 1), e.Description())
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

		kapi := etcd.NewKeysAPI(client)
		if _, err = kapi.Delete(context.Background(), strings.TrimRight(*dir, "/"), &etcd.DeleteOptions{Recursive: true}); err != nil {
			log.Fatalf(err.Error())
		}
		log.Printf("Removed path: %s", strings.TrimRight(*dir, "/"))
	}

	// Create dir.
	kapi := etcd.NewKeysAPI(client)
	if _, err := kapi.Set(context.TODO(), *dir, "", &etcd.SetOptions{Dir: true}); err != nil {
		log.Printf(err.Error())

		// Should prob. check that we're actually dealing with an existing key and not something else...
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
