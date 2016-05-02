package command

import (
	"strings"

	"github.com/codegangsta/cli"
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
		Action: func(c *cli.Context) error {
			exportCommandFunc(c)
			return nil
		},
	}
}

// exportCommandFunc exports data as either JSON, YAML or TOML.
func exportCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		fatal("You need to specify directory")
	}
	dir := c.Args()[0]

	// Remove trailing slash.
	if dir != "/" {
		dir = strings.TrimRight(dir, "/")
	}
	infof("Using dir: %s", dir)

	// Load configuration file.
	e := loadConfig(c)

	// New dir API.
	ki := newKeyAPI(e)

	sort := c.Bool("sort")

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		fatal(err.Error())
	}

	exportFunc(dir, sort, c.String("output"), f, c, ki)
}

// exportCommandFunc exports data as either JSON, YAML or TOML.
func exportFunc(dir string, sort bool, file string, f iodatafmt.DataFmt, c *cli.Context, ki client.KeysAPI) {
	ctx, cancel := contextWithCommandTimeout(c)
	resp, err := ki.Get(ctx, dir, &client.GetOptions{Sort: sort, Recursive: true})
	cancel()
	if err != nil {
		fatal(err.Error())
	}

	// Export and write output.
	m := etcdmap.Map(resp.Node)
	if file != "" {
		iodatafmt.Write(file, m, f)
	} else {
		iodatafmt.Print(m, f)
	}
}
