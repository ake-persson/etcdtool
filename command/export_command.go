package command

import (
	"strings"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
)

// NewExportCommand returns data from export.
func NewExportCommand() cli.Command {
	return cli.Command{
		Name:  "export",
		Usage: "export a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort, s", Usage: "returns result in sorted order"},
			cli.StringFlag{Name: "format, f", EnvVar: "ETCDTOOL_FORMAT", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "output, o", Value: "", Usage: "Output file"},
		},
		Action: func(c *cli.Context) {
			exportCommandFunc(c, mustNewKeyAPI(c))
		},
	}
}

// exportCommandFunc exports data as either JSON, YAML or TOML.
func exportCommandFunc(c *cli.Context, ki client.KeysAPI) {
	key := "/"
	if len(c.Args()) != 0 {
		key = strings.TrimRight(c.Args()[0], "/") + "/"
	}

	sort := c.Bool("sort")

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		handleError(ExitServerError, err)
	}

	exportFunc(key, sort, c.String("output"), f, c, ki)
}

// exportCommandFunc exports data as either JSON, YAML or TOML.
func exportFunc(key string, sort bool, file string, f iodatafmt.DataFmt, c *cli.Context, ki client.KeysAPI) {
	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Get(ctx, key, &client.GetOptions{Sort: sort, Recursive: true})
	cancel()
	if err != nil {
		handleError(ExitServerError, err)
	}

	// Export and write output.
	m := etcdmap.Map(resp.Node)
	if file != "" {
		iodatafmt.Write(file, m, f)
	} else {
		iodatafmt.Print(m, f)
	}
}
