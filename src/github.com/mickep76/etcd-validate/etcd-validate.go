package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	jsonschema "github.com/xeipuuv/gojsonschema"

	"github.com/mickep76/common"
)

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

	if *schema == "" {
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
			log.Fatalf("Couldn't determine schema to use for directory: %s", *dir)
		}
	}

	// Get JSON Schema.
	kapi := etcd.NewKeysAPI(client)
	res, err := kapi.Get(context.Background(), *schema, &etcd.GetOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	schemaLoader := jsonschema.NewStringLoader(res.Node.Value)

	// Get etcd dir.
	res2, err2 := kapi.Get(context.Background(), *dir, &etcd.GetOptions{Recursive: true})
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
			fmt.Printf("%s: %s\n", strings.Replace(e.Context().String("/"), "(root)", *dir, 1), e.Description())
		}
	}
}
