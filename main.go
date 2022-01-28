package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mickep76/etcdtool/command"
)

type Config struct {
	peers string
}

func readStandardConfig() Config {
	var ret Config
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	cfg := path.Join(dir, ".etcdtool.conf")

	file, err := os.Open(cfg)
	if err != nil {
		return ret
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		keyvals := strings.Split(line, " ")
		if len(keyvals) != 2 {
			log.Fatalf("Bad line in config file %s, line: %s", cfg, line)
		}
		if keyvals[0] == "peers" {
			ret.peers = keyvals[1]
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ret
}

func main() {
	cfg := readStandardConfig()
	defaultPeers := "http://127.0.0.1:4001,http://127.0.0.1:2379"
	if len(cfg.peers) > 0 {
		defaultPeers = cfg.peers
	}

	app := cli.NewApp()
	app.Name = "etcdtool"
	app.Version = Version
	app.Usage = "Command line tool for etcd to import, export, edit or validate data in either JSON, YAML or TOML format."
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", EnvVar: "ETCDTOOL_CONFIG", Usage: "Configuration file"},
		cli.BoolFlag{Name: "debug, d", Usage: "Debug"},
		cli.StringFlag{Name: "peers, p", Value: defaultPeers, EnvVar: "ETCDTOOL_PEERS", Usage: "Comma-delimited list of hosts in the cluster"},
		cli.StringFlag{Name: "cert", Value: "", EnvVar: "ETCDTOOL_CERT", Usage: "Identify HTTPS client using this SSL certificate file"},
		cli.StringFlag{Name: "key", Value: "", EnvVar: "ETCDTOOL_KEY", Usage: "Identify HTTPS client using this SSL key file"},
		cli.StringFlag{Name: "ca", Value: "", EnvVar: "ETCDTOOL_CA", Usage: "Verify certificates of HTTPS-enabled servers using this CA bundle"},
		cli.StringFlag{Name: "user, u", Value: "", Usage: "User"},
		cli.StringFlag{Name: "password-file, F", Value: "", Usage: "File path to the user's password"},
		cli.DurationFlag{Name: "timeout, t", Value: time.Second, Usage: "Connection timeout"},
		cli.DurationFlag{Name: "command-timeout, T", Value: 5 * time.Second, Usage: "Command timeout"},
	}
	app.Commands = []cli.Command{
		command.NewImportCommand(),
		command.NewExportCommand(),
		command.NewLsCommand(),
		command.NewEditCommand(),
		command.NewValidateCommand(),
		command.NewTreeCommand(),
		command.NewPrintConfigCommand(),
	}

	app.Run(os.Args)
}
