// Package etcdmap provides methods for interacting with Etcd using struct, map or JSON.
package etcdmap

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"encoding/json"
	"github.com/coreos/go-etcd/etcd"
)

// Struct returns a struct from a Etcd directory.
// !!! This is not supported for nested struct yet.
func Struct(root *etcd.Node, s interface{}) error {
	// Convert Etcd node to map[string]interface{}
	m := Map(root)

	// Yes this is a hack, so what it works.
	// Marshal map[string]interface{} to JSON.
	j, err := json.Marshal(&m)
	if err != nil {
		return err
	}

	// Yes this is a hack, so what it works.
	// Unmarshal JSON to struct.
	if err := json.Unmarshal(j, &s); err != nil {
		return err
	}

	return nil
}

// JSON returns an Etcd directory as JSON []byte.
func JSON(root *etcd.Node) ([]byte, error) {
	j, err := json.Marshal(Map(root))
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// JSONIndent returns an Etcd directory as indented JSON []byte.
func JSONIndent(root *etcd.Node, indent string) ([]byte, error) {
	j, err := json.MarshalIndent(Map(root), "", indent)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// Map returns a map[string]interface{} from a Etcd directory.
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

// CreateStruct creates a Etcd directory based on a struct.
func CreateStruct(client *etcd.Client, dir string, s interface{}) error {
	// Yes this is a hack, so what it works.
	// Marshal struct to JSON
	j, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	// Yes this is a hack, so what it works.
	// Unmarshal JSON to map[string]interface{}
	m := make(map[string]interface{})
	if err := json.Unmarshal(j, &m); err != nil {
		return err
	}

	return CreateMap(client, dir, m)
}

// CreateJSON creates a Etcd directory based on JSON byte[].
func CreateJSON(client *etcd.Client, dir string, j []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(j, &m); err != nil {
		return err
	}

	return CreateMap(client, dir, m)
}

// CreateMap creates a Etcd directory based on map[string]interface{}.
func CreateMap(client *etcd.Client, dir string, d map[string]interface{}) error {
	for k, v := range d {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			if _, err := client.CreateDir(dir+"/"+k, 0); err != nil {
				return err
			}
			CreateMap(client, dir+"/"+k, v.(map[string]interface{}))
		} else if reflect.ValueOf(v).Kind() == reflect.Slice {
			CreateMapSlice(client, dir+"/"+k, v.([]interface{}))
		} else if reflect.ValueOf(v).Kind() == reflect.String {
			if _, err := client.Set(dir+"/"+k, v.(string), 0); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported type: %s for key: %s", reflect.ValueOf(v).Kind(), k)
		}
	}

	return nil
}

// CreateMapSlice creates a Etcd directory based on []interface{}.
func CreateMapSlice(client *etcd.Client, dir string, d []interface{}) error {
	for i, v := range d {
		istr := strconv.Itoa(i)
		if reflect.ValueOf(v).Kind() == reflect.Map {
			if _, err := client.CreateDir(dir+"/"+istr, 0); err != nil {
				return err
			}
			CreateMap(client, dir+"/"+istr, v.(map[string]interface{}))
		} else if reflect.ValueOf(v).Kind() == reflect.Slice {
			CreateMapSlice(client, dir+"/"+istr, v.([]interface{}))
		} else if reflect.ValueOf(v).Kind() == reflect.String {
			if _, err := client.Set(dir+"/"+istr, v.(string), 0); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported type: %s for key: %d", reflect.ValueOf(v).Kind(), i)
		}
	}

	return nil
}
