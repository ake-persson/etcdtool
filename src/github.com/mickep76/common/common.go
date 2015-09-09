package common

import (
	"os"
	"strings"
)

// GetEnv variable for Etcd connection.
func GetEnv() []string {
	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		if a[0] == "ETCD_CONN" {
			return []string{a[1]}
		}
	}

	return []string{}
}
