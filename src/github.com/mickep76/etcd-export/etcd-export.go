package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"

	"github.com/mickep76/common"
)

func main() {
	// Options.
	version := flag.Bool("version", false, "Version")
	peers := flag.String("peers", common.GetEnv(), "Comma separated list of etcd nodes")
	//	dir := flag.String("dir", "/", "etcd directory")
	format := flag.String("format", "JSON", "Data serialization format YAML, TOML or JSON")
	output := flag.String("output", "", "Output file")
	flag.Parse()

	var dir string
	if len(flag.Args()) < 1 {
		log.Fatal("You need to specify dir.")
	}
	dir = flag.Args()[0]

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

	// Export data.
	kapi := etcd.NewKeysAPI(client)
	res, err3 := kapi.Get(context.Background(), dir, &etcd.GetOptions{Recursive: true})
	if err3 != nil {
		log.Fatal(err3.Error())
	}
	m := etcdmap.Map(res.Node)

	// Write output.
	if *output != "" {
		iodatafmt.Write(*output, m, f)
	} else {
		iodatafmt.Print(m, f)
	}
}
