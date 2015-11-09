package main

import (
	"os"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/mickep76/etcdfmt/command"
)

func main() {
	app := cli.NewApp()
	app.Name = "etcdfmt"
	app.Version = Version
	app.Usage = "Command line tool for etcd to import, export, edit or validate data in either JSON, YAML or TOML format."
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "output cURL commands which can be used to reproduce the request"},
		cli.StringFlag{Name: "peers, C", Value: "", Usage: "a comma-delimited list of machine addresses in the cluster (default: \"http://127.0.0.1:4001,http://127.0.0.1:2379\")"},
		cli.StringFlag{Name: "endpoint", Value: "", Usage: "a comma-delimited list of machine addresses in the cluster (default: \"http://127.0.0.1:4001,http://127.0.0.1:2379\")"},
		cli.StringFlag{Name: "cert-file", Value: "", Usage: "identify HTTPS client using this SSL certificate file"},
		cli.StringFlag{Name: "key-file", Value: "", Usage: "identify HTTPS client using this SSL key file"},
		cli.StringFlag{Name: "ca-file", Value: "", Usage: "verify certificates of HTTPS-enabled servers using this CA bundle"},
		cli.StringFlag{Name: "username, u", Value: "", Usage: "provide username[:password] and prompt if password is not supplied."},
		cli.DurationFlag{Name: "timeout", Value: time.Second, Usage: "connection timeout per request"},
		cli.DurationFlag{Name: "total-timeout", Value: 5 * time.Second, Usage: "timeout for the command execution (except watch)"},
	}
	app.Commands = []cli.Command{
		command.NewImportCommand(),
		command.NewExportCommand(),
		command.NewEditCommand(),
	}

	app.Run(os.Args)
}
