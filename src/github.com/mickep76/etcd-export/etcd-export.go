package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/go-etcd/etcd"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

// EtcdMap creates a nested data structure from a Etcd node.
func EtcdMap(root *etcd.Node) map[string]interface{} {
	v := make(map[string]interface{})

	for _, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]
		if n.Dir {
			v[k] = make(map[string]interface{})
			v[k] = EtcdMap(n)
		} else {
			v[k] = n.Value
		}
	}
	return v
}

func main() {
	// Set log options.
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

	// Options.
	var opts struct {
		Verbose  bool    `short:"v" long:"verbose" description:"Verbose"`
		Version  bool    `long:"version" description:"Version"`
		Format   string  `short:"f" long:"format" description:"Data serialization format YAML, TOML or JSON" default:"YAML"`
		Output   *string `short:"o" long:"output" description:"Output file (STDOUT)"`
		EtcdHost *string `short:"H" long:"etcd-host" description:"Etcd Host"`
		EtcdPort int     `short:"p" long:"etcd-port" description:"Etcd Port" default:"2379"`
		EtcdDir  string  `short:"d" long:"etcd-dir" description:"Etcd Dir" default:"/"`
	}

	// Parse options.
	if _, err := flags.Parse(&opts); err != nil {
		ferr := err.(*flags.Error)
		if ferr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Fatal(err.Error())
		}
	}

	// Print version.
	if opts.Version {
		fmt.Printf("tf %s\n", Version)
		os.Exit(0)
	}

	// Set verbose.
	if opts.Verbose {
		log.SetLevel(log.InfoLevel)
	}

	// Validate input.
	if opts.EtcdHost == nil {
		log.Fatalf("You need to specify Etcd host.")
	}

	// Get Etcd input.
	node := []string{fmt.Sprintf("http://%v:%v", *opts.EtcdHost, opts.EtcdPort)}
	client := etcd.NewClient(node)
	res, err := client.Get(opts.EtcdDir, true, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	data := EtcdMap(res.Node)

	switch strings.ToUpper(opts.Format) {
	case "YAML":
		s, err := yaml.Marshal(&data)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(s))
	case "TOML":
		s := new(bytes.Buffer)
		err := toml.NewEncoder(s).Encode(&data)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(s.String()))
	case "JSON":
		s, err := json.MarshalIndent(&data, "", "    ")
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(s), "\n")
	default:
		log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
	}
}
