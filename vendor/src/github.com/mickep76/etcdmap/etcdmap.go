// Package etcdmap provides methods for interacting with Etcd using struct, map or JSON.
package etcdmap

import (
	"fmt"
	"reflect"
	"strings"

	"encoding/json"
	//	"github.com/coreos/go-etcd/etcd"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
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

// Create Etcd directory structure from a map, slice or struct.
func Create(client *etcd.Client, path string, val reflect.Value) error {
	switch val.Kind() {
	case reflect.Ptr:
		orig := val.Elem()
		if !orig.IsValid() {
			return nil
		}
		if err := Create(client, path, orig); err != nil {
			return err
		}
	case reflect.Interface:
		orig := val.Elem()
		if err := Create(client, path, orig); err != nil {
			return err
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			t := val.Type().Field(i)
			k := t.Tag.Get("etcd")
			if err := Create(client, path+"/"+k, val.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			if err := Create(client, path+"/"+k.String(), v); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			Create(client, fmt.Sprintf("%s/%d", path, i), val.Index(i))
		}
	case reflect.String:
		kapi := etcd.NewKeysAPI(*client)
		_, err := kapi.Set(context.Background(), path, val.String(), nil)
		if err != nil {
			return err
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		kapi := etcd.NewKeysAPI(*client)
		_, err := kapi.Set(context.Background(), path, fmt.Sprintf("%v", val.Interface()), nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported type: %s for path: %s", val.Kind(), path)
	}

	return nil
}

// CreateJSON Etcd directory structure from JSON.
func CreateJSON(client *etcd.Client, dir string, j []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(j, &m); err != nil {
		return err
	}

	return Create(client, dir, reflect.ValueOf(m))
}
