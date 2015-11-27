package command

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"
	"golang.org/x/net/context"
)

func contextWithCommandTimeout(c *cli.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.GlobalDuration("command-timeout"))
}

func newTransport(c *cli.Context) *http.Transport {
	tls := transport.TLSInfo{
		CAFile:   c.GlobalString("ca"),
		CertFile: c.GlobalString("cert"),
		KeyFile:  c.GlobalString("key"),
	}

	timeout := 30 * time.Second
	tr, err := transport.NewTransport(tls, timeout)
	if err != nil {
		log.Fatal(err.Error())
	}

	return tr
}

func newClient(c *cli.Context) client.Client {
	cfg := client.Config{
		Transport:               newTransport(c),
		Endpoints:               strings.Split(c.GlobalString("peers"), ","),
		HeaderTimeoutPerRequest: c.GlobalDuration("timeout"),
	}

	/*
	   uFlag := c.GlobalString("username")
	   if uFlag != "" {
	       username, password, err := getUsernamePasswordFromFlag(uFlag)
	       if err != nil {
	           return nil, err
	       }
	       cfg.Username = username
	       cfg.Password = password
	   }
	*/
	cl, err := client.New(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	return cl
}

func newKeyAPI(c *cli.Context) client.KeysAPI {
	return client.NewKeysAPI(newClient(c))
}
