package etcdmap

import (
	"reflect"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"
)

// Map creates a map[string]interface{} from a Etcd directory.
func Map(root *etcd.Node) map[string]interface{} {
	v := make(map[string]interface{})

	for _, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]
		if n.Dir {
			v[k] = make(map[string]interface{})
			v[k] = Map(n)
		} else {
			v[k] = n.Value
		}
	}
	return v
}

// CreateMap create Etcd directory structure using a map[string]interface{}.
func CreateMap(client *etcd.Client, dir string, d map[string]interface{}) error {
	for k, v := range d {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			if _, err := client.CreateDir(dir+"/"+k, 0); err != nil {
				return err
			}
			CreateMap(client, dir+"/"+k, v.(map[string]interface{}))
		} else {
			if _, err := client.Set(dir+"/"+k, v.(string), 0); err != nil {
				return err
			}
		}
	}

	return nil
}
