package command

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bgentry/speakeasy"
	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"
	"golang.org/x/net/context"
)

func contextWithCommandTimeout(c *cli.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.GlobalDuration("command-timeout"))
}

func newTransport(e Etcdtool) *http.Transport {
	tls := transport.TLSInfo{
		CAFile:   e.CA,
		CertFile: e.Cert,
		KeyFile:  e.Key,
	}

	timeout := 30 * time.Second
	tr, err := transport.NewTransport(tls, timeout)
	if err != nil {
		log.Fatal(err.Error())
	}

	return tr
}

func newClient(e Etcdtool) client.Client {
	cfg := client.Config{
		Transport:               newTransport(e),
		Endpoints:               strings.Split(e.Peers, ","),
		HeaderTimeoutPerRequest: e.Timeout,
	}

	if e.User != "" {
		cfg.Username = e.User
		var err error
		cfg.Password, err = speakeasy.Ask("Password: ")
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	cl, err := client.New(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	return cl
}

func newKeyAPI(e Etcdtool) client.KeysAPI {
	return client.NewKeysAPI(newClient(e))
}
