package command

import (
	"log"
	"os"
	"os/user"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mickep76/iodatafmt"
)

type etcdtool struct {
	peers          string        `json:"peers,omitempty" yaml:"peers,omitempty" toml:"peers,omitempty"`
	cert           string        `json:"cert,omitempty" yaml:"cert,omitempty" toml:"cert,omitempty"`
	key            string        `json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	ca             string        `json:"ca,omitempty" yaml:"ca,omitempty" toml:"peers,omitempty"`
	user           string        `json:"user,omitempty" yaml:"user,omitempty" toml:"user,omitempty"`
	timeout        time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" toml:"timeout,omitempty"`
	commandTimeout time.Duration `json:"commandTimeout,omitempty" yaml:"commandTimeout,omitempty" toml:"commandTimoeut,omitempty"`
	routes         []route       `json:"routes" yaml:"routes" toml:"routes"`
}

type route struct {
	regexp string `json:"regexp" yaml:"regexp" toml:"regexp"`
	schema string `json:"schema" yaml:"schema" toml:"schema"`
}

func LoadConfig(c *cli.Context) etcdtool {
	if c.GlobalString("config") != "" {
		infof(c, "Using config file: %s", c.GlobalString("config"))
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

	e := etcdtool{}
	for _, fn := range cfgs {
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			continue
		}
		infof(c, "Using config file: %s", fn)
		f, err := iodatafmt.FileFormat(fn)
		if err != nil {
			log.Fatal(err.Error())
		}
		if err := iodatafmt.LoadPtr(&e, fn, f); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Override with arguments or env. variables.
	if c.GlobalString("peers") != "" {
		e.peers = c.GlobalString("peers")
	}

	if c.GlobalString("cert") != "" {
		e.cert = c.GlobalString("cert")
	}

	if c.GlobalString("key") != "" {
		e.key = c.GlobalString("key")
	}

	if c.GlobalString("ca") != "" {
		e.ca = c.GlobalString("ca")
	}

	if c.GlobalString("user") != "" {
		e.user = c.GlobalString("user")
	}

	if c.IsSet("timeout") {
		e.timeout = c.GlobalDuration("timeout")
	}

	if c.IsSet("command-timeout") {
		e.commandTimeout = c.GlobalDuration("command-timeout")
	}

	return e
}
