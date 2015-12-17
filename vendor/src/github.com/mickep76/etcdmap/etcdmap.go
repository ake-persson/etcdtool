// Package etcdmap provides methods for interacting with etcd using struct, map or JSON.
package etcdmap

import (
	"encoding/json"
	"fmt"
	"log"
	pathx "path"
	"reflect"
	"strconv"
	"strings"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

func strToInt(v string, k reflect.Kind) (reflect.Value, error) {
	switch k {
	case reflect.Int:
		i, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int(i)), nil
	case reflect.Int8:
		i, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int8(i)), nil
	case reflect.Int16:
		i, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int16(i)), nil
	case reflect.Int32:
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int32(i)), nil
	case reflect.Int64:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(int64(i)), nil
	}

	return reflect.Value{}, fmt.Errorf("unsupported type: %s", k)
}

func strToUint(v string, k reflect.Kind) (reflect.Value, error) {
	switch k {
	case reflect.Uint:
		i, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint(i)), nil
	case reflect.Uint8:
		i, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint8(i)), nil
	case reflect.Uint16:
		i, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint16(i)), nil
	case reflect.Uint32:
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint32(i)), nil
	case reflect.Uint64:
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(uint64(i)), nil
	}

	return reflect.Value{}, fmt.Errorf("unsupported type: %s", k)
}

func strToFloat(v string, k reflect.Kind) (reflect.Value, error) {
	switch k {
	case reflect.Float32:
		i, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(float32(i)), nil
	case reflect.Float64:
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(float64(i)), nil
	}

	return reflect.Value{}, fmt.Errorf("unsupported type: %s", k)
}

func Struct(root *client.Node, val reflect.Value) error {
	m := Map(root)
	s := val.Elem()

	for i := 0; i < s.NumField(); i++ {
		t := s.Type().Field(i)
		k := t.Tag.Get("etcd")
		if v, ok := m[k]; ok {
			f := s.Field(i)
			switch f.Type().Kind() {
			case reflect.String:
				f.Set(reflect.ValueOf(v))
			case reflect.Bool:
				i, err := strconv.ParseBool(v.(string))
				if err != nil {
					return err
				}
				f.SetBool(i)
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strToInt(v.(string), f.Type().Kind())
				if err != nil {
					return err
				}
				f.Set(i)
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err := strToUint(v.(string), f.Type().Kind())
				if err != nil {
					return err
				}
				f.Set(i)
			case reflect.Float32, reflect.Float64:
				i, err := strToFloat(v.(string), f.Type().Kind())
				if err != nil {
					return err
				}
				f.Set(i)
			default:
				return fmt.Errorf("unsupported type: %s for key: %s", f.Type().Kind(), k)
			}
		}
	}
	return nil
}

// JSON returns an etcd directory as JSON []byte.
func JSON(root *client.Node) ([]byte, error) {
	j, err := json.Marshal(Map(root))
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// JSON returns an etcd directory as JSON []byte.
func ArrayJSON(root *client.Node) ([]byte, error) {
	j, err := json.Marshal(Array(root))
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// JSONIndent returns an etcd directory as indented JSON []byte.
func JSONIndent(root *client.Node, indent string) ([]byte, error) {
	j, err := json.MarshalIndent(Map(root), "", indent)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// JSONIndent returns an etcd directory as indented JSON []byte.
func ArrayJSONIndent(root *client.Node, indent string) ([]byte, error) {
	j, err := json.MarshalIndent(Array(root), "", indent)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

// Map returns a map[string]interface{} from a etcd directory.
func Map(root *client.Node) map[string]interface{} {
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

// Array returns a []interface{} including the directory name inside each entry from a etcd directory.
func Array(root *client.Node) []interface{} {
	v := []interface{}{}

	for _, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]
		if n.Dir {
			m := make(map[string]interface{})
			m = Map(n)
			m["dir_name"] = k
		}
	}
	return v
}

// Create etcd directory structure from a map, slice or struct.
func Create(kapi client.KeysAPI, path string, val reflect.Value) error {
	switch val.Kind() {
	case reflect.Ptr:
		orig := val.Elem()
		if !orig.IsValid() {
			return nil
		}
		if err := Create(kapi, path, orig); err != nil {
			return err
		}
	case reflect.Interface:
		orig := val.Elem()
		if err := Create(kapi, path, orig); err != nil {
			return err
		}
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			t := val.Type().Field(i)
			k := t.Tag.Get("etcd")
			if err := Create(kapi, path+"/"+k, val.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		if strings.HasPrefix(pathx.Base(path), "_") {
			log.Printf("create hidden directory in etcd: %s", path)
		}
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			if err := Create(kapi, path+"/"+k.String(), v); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			Create(kapi, fmt.Sprintf("%s/%d", path, i), val.Index(i))
		}
	case reflect.String:
		if strings.HasPrefix(pathx.Base(path), "_") {
			log.Printf("set hidden key in etcd: %s", path)
		}
		_, err := kapi.Set(context.TODO(), path, val.String(), nil)
		if err != nil {
			return err
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		if strings.HasPrefix(pathx.Base(path), "_") {
			log.Printf("set hidden key in etcd: %s", path)
		}
		_, err := kapi.Set(context.TODO(), path, fmt.Sprintf("%v", val.Interface()), nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported type: %s for path: %s", val.Kind(), path)
	}

	return nil
}

// CreateJSON etcd directory structure from JSON.
func CreateJSON(kapi client.KeysAPI, dir string, j []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(j, &m); err != nil {
		return err
	}

	return Create(kapi, dir, reflect.ValueOf(m))
}
