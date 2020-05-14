package command

import (
	"os"
	"os/user"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mickep76/iodatafmt"
)

// Etcdtool configuration struct.
type Etcdtool struct {
	Peers            string        `json:"peers,omitempty" yaml:"peers,omitempty" toml:"peers,omitempty"`
	Cert             string        `json:"cert,omitempty" yaml:"cert,omitempty" toml:"cert,omitempty"`
	Key              string        `json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	CA               string        `json:"ca,omitempty" yaml:"ca,omitempty" toml:"peers,omitempty"`
	User             string        `json:"user,omitempty" yaml:"user,omitempty" toml:"user,omitempty"`
	Timeout          time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" toml:"timeout,omitempty"`
	CommandTimeout   time.Duration `json:"commandTimeout,omitempty" yaml:"commandTimeout,omitempty" toml:"commandTimeout,omitempty"`
	Routes           []Route       `json:"routes" yaml:"routes" toml:"routes"`
	PasswordFilePath string
}

// Route configuration struct.
type Route struct {
	Regexp string `json:"regexp" yaml:"regexp" toml:"regexp"`
	Schema string `json:"schema" yaml:"schema" toml:"schema"`
}

func loadConfig(c *cli.Context) Etcdtool {
	// Enable debug
	if c.GlobalBool("debug") {
		debug = true
	}

	// Default path for config file.
	u, _ := user.Current()
	cfgs := []string{
		"/etcd/etcdtool.json",
		"/etcd/etcdtool.yaml",
		"/etcd/etcdtool.toml",
		u.HomeDir + "/.etcdtool.json",
		u.HomeDir + "/.etcdtool.yaml",
		u.HomeDir + "/.etcdtool.toml",
	}

	// Check if we have an arg. for config file and that it exist's.
	if c.GlobalString("config") != "" {
		if _, err := os.Stat(c.GlobalString("config")); os.IsNotExist(err) {
			fatalf("Config file doesn't exist: %s", c.GlobalString("config"))
		}
		cfgs = append([]string{c.GlobalString("config")}, cfgs...)
	}

	// Check if config file exists and load it.
	e := Etcdtool{}
	for _, fn := range cfgs {
		if _, err := os.Stat(fn); os.IsNotExist(err) {
			continue
		}
		infof("Using config file: %s", fn)
		f, err := iodatafmt.FileFormat(fn)
		if err != nil {
			fatal(err.Error())
		}
		if err := iodatafmt.LoadPtr(&e, fn, f); err != nil {
			fatal(err.Error())
		}
	}

	// Override with arguments or env. variables.
	if c.GlobalString("peers") != "" {
		e.Peers = c.GlobalString("peers")
	}

	if c.GlobalString("cert") != "" {
		e.Cert = c.GlobalString("cert")
	}

	if c.GlobalString("key") != "" {
		e.Key = c.GlobalString("key")
	}

	if c.GlobalString("ca") != "" {
		e.CA = c.GlobalString("ca")
	}

	if c.GlobalString("user") != "" {
		e.User = c.GlobalString("user")
	}

	if c.GlobalDuration("timeout") != 0 {
		e.Timeout = c.GlobalDuration("timeout")
	}

	if c.GlobalDuration("command-timeout") != 0 {
		e.CommandTimeout = c.GlobalDuration("command-timeout")
	}

	// Add password file path if set
	if c.GlobalString("password-file") != "" {
		e.PasswordFilePath = c.GlobalString("password-file")
	}

	return e
}
