package common

import (
	"os"
	"strings"
)

// GetEnv variable for Etcd connection.
func GetEnv() string {
	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		if a[0] == "ETCD_PEERS" {
			return a[1]
		} else if a[0] == "ETCDCTL_PEERS" {
			return a[1]
		}
	}

	return "127.0.0.1:4001,127.0.0.1:2379"
}
