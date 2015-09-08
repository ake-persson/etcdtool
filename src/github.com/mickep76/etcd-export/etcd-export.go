package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/go-etcd/etcd"
	flags "github.com/jessevdk/go-flags"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
)

func getEnv() []string {
	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		if a[0] == "ETCD_CONN" {
			return []string{a[1]}
		}
	}

	return []string{}
}

func main() {
	// Set log options.
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

	// Get connection env variable.
	conn := getEnv()

	// Options.
	var opts struct {
		Verbose  bool    `short:"v" long:"verbose" description:"Verbose"`
		Version  bool    `long:"version" description:"Version"`
		Format   string  `short:"f" long:"format" description:"Data serialization format YAML, TOML or JSON" default:"JSON"`
		Output   *string `short:"o" long:"output" description:"Output file (STDOUT)"`
		EtcdNode *string `short:"n" long:"etcd-node" description:"Etcd Node"`
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
		fmt.Printf("etcd-export %s\n", Version)
		os.Exit(0)
	}

	// Set verbose.
	if opts.Verbose {
		log.SetLevel(log.InfoLevel)
	}

	// Validate input.
	if len(conn) < 1 && opts.EtcdNode == nil {
		log.Fatalf("You need to specify Etcd host.")
	}

	// Get data format.
	f, err := iodatafmt.Format(opts.Format)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Setup Etcd client.
	if opts.EtcdNode != nil {
		conn = []string{fmt.Sprintf("http://%v:%v", *opts.EtcdNode, opts.EtcdPort)}
	}
	client := etcd.NewClient(conn)

	// Export data.
	res, err := client.Get(opts.EtcdDir, true, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := etcdmap.Map(res.Node)

	// Write output.
	if opts.Output != nil {
		iodatafmt.Write(*opts.Output, d, f)
	} else {
		iodatafmt.Print(d, f)
	}
}
