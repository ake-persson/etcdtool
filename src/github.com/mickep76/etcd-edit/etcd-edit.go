package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"

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

	//	fmt.Println(*editor, *tmpFile)

	cmd := exec.Command(*editor, *tmpFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err3 := cmd.Run()
	if err3 != nil {
		log.Fatalf(err3.Error())
	}

	/*

		cmd := exec.Command(*editor, *tmpFile)
		err3 := cmd.Start()
		if err3 != nil {
			log.Fatal(err3)
		}
		log.Printf("Waiting for command to finish...")
		err4 := cmd.Wait()
		if err4 != nil {
			log.Fatal(err4)
		}
	*/

	/*
		err3 := exec.Command(*editor, *tmpFile).Run()
		if err3 != nil {
			panic(err3)
			//		log.Fatalf(err3.Error())
		}
	*/
}
