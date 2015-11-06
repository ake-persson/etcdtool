package command

import (
	"log"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
	"github.com/mickep76/iodatafmt"
)

func NewExportCommand() cli.Command {
	return cli.Command{
		Name:  "export",
		Usage: "export a directory",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort", Usage: "returns result in sorted order"},
			cli.StringFlag{Name: "format", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
			cli.StringFlag{Name: "output", Value: "", Usage: "Output file"},
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
		key = c.Args()[0]
	}

	sort := c.Bool("sort")

	ctx, cancel := contextWithTotalTimeout(c)
	resp, err := ki.Get(ctx, key, &client.GetOptions{Sort: sort, Recursive: true})
	cancel()
	if err != nil {
		handleError(ExitServerError, err)
	}

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Export data.
	m := etcdmap.Map(resp.Node)
	iodatafmt.Print(m, f)
}
