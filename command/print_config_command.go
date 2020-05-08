package command

import (
	"github.com/codegangsta/cli"
	"github.com/mickep76/iodatafmt"
)

// NewPrintConfigCommand print configuration.
func NewPrintConfigCommand() cli.Command {
	return cli.Command{
		Name:  "print-config",
		Usage: "Print configuration",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "sort", Usage: "returns result in sorted order"},
			cli.StringFlag{Name: "format, f", EnvVar: "ETCDTOOL_FORMAT", Value: "JSON", Usage: "Data serialization format YAML, TOML or JSON"},
		},
		Action: func(c *cli.Context) error {
			printConfigCommandFunc(c)
			return nil
		},
	}
}

func printConfigCommandFunc(c *cli.Context) {
	// Load configuration file.
	e := loadConfig(c)

	// Get data format.
	f, err := iodatafmt.Format(c.String("format"))
	if err != nil {
		fatal(err.Error())
	}

	iodatafmt.Print(e, f)
}
