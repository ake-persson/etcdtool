# EtcdMap

Go package provides methods for interacting with Etcd using struct, map or JSON.

[![GoDoc](https://godoc.org/github.com/mickep76/etcdmap?status.svg)](https://godoc.org/github.com/mickep76/etcdmap)

# Documentation


# etcdmap
    import "github.com/mickep76/etcdmap"

Package etcdmap provides methods for interacting with etcd using struct, map or JSON.






## func Array
``` go
func Array(root *client.Node) []interface{}
```
Array returns a []interface{} including the directory name inside each entry from a etcd directory.


## func ArrayJSON
``` go
func ArrayJSON(root *client.Node) ([]byte, error)
```
JSON returns an etcd directory as JSON []byte.


## func ArrayJSONIndent
``` go
func ArrayJSONIndent(root *client.Node, indent string) ([]byte, error)
```
JSONIndent returns an etcd directory as indented JSON []byte.


## func Create
``` go
func Create(kapi client.KeysAPI, path string, val reflect.Value) error
```
Create etcd directory structure from a map, slice or struct.


## func CreateJSON
``` go
func CreateJSON(kapi client.KeysAPI, dir string, j []byte) error
```
CreateJSON etcd directory structure from JSON.


## func JSON
``` go
func JSON(root *client.Node) ([]byte, error)
```
JSON returns an etcd directory as JSON []byte.


## func JSONIndent
``` go
func JSONIndent(root *client.Node, indent string) ([]byte, error)
```
JSONIndent returns an etcd directory as indented JSON []byte.


## func Map
``` go
func Map(root *client.Node) map[string]interface{}
```
Map returns a map[string]interface{} from a etcd directory.


## func Struct
``` go
func Struct(root *client.Node, val reflect.Value) error
```








- - -
