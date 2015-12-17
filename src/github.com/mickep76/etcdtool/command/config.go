package command

import (
	"log"
	"os"
	"os/user"

	"github.com/codegangsta/cli"
	"github.com/mickep76/iodatafmt"
)

type Etcdtool struct {
	Peers  string  `json:"peers" yaml:"peers" toml:"peers"`
	Routes []Route `json:"routes" yaml:"peers" toml:"routes"`
}

type Route struct {
	Regexp string `json:"regexp" yaml:"regexp" toml:"regexp"`
	Schema string `json:"schema" yaml:"schema" toml:"schema"`
}

func LoadConfig(c *cli.Context) Etcdtool {
	if c.GlobalString("config") != "" {
		Infof(c, "Using config file: %s", c.GlobalString("config"))
		if _, err := os.Stat(c.GlobalString("config")); os.IsNotExist(err) {
			log.Fatalf("Config file doesn't exist: %s", c.GlobalString("config"))
		}
	}

	u, _ := user.Current()
	cfgs := []string{
		"/etcd/etcdtool.json",
		"/etcd/etcdtool.yaml",
		"/etcd/etcdtool.toml",
		u.HomeDir + "/.etcdtool.json",
		u.HomeDir + "/.etcdtool.yaml",
		u.HomeDir + "/.etcdtool.toml",
	}

	s := Etcdtool{}
	for _, fn := range cfgs {
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			continue
		}
		Infof(c, "Using config file: %s", fn)
		f, err := iodatafmt.FileFormat(fn)
		if err != nil {
			log.Fatal(err.Error())
		}
		if err := iodatafmt.LoadPtr(&s, fn, f); err != nil {
			log.Fatal(err.Error())
		}
	}

	return s
}
