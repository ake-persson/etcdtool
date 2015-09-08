# EtcdMap

Go package provides methods for interacting with Etcd using struct, map or JSON.

[![GoDoc](https://godoc.org/github.com/mickep76/etcdmap?status.svg)](https://godoc.org/github.com/mickep76/etcdmap)

# Documentation


# etcdmap
    import "github.com/mickep76/etcdmap"

Package etcdmap provides methods for interacting with Etcd using struct, map or JSON.






## func CreateJSON
``` go
func CreateJSON(client *etcd.Client, dir string, j []byte) error
```
CreateJSON creates a Etcd directory based on JSON byte[].


## func CreateMap
``` go
func CreateMap(client *etcd.Client, dir string, d map[string]interface{}) error
```
CreateMap creates a Etcd directory based on map[string]interface{}.


## func CreateMapSlice
``` go
func CreateMapSlice(client *etcd.Client, dir string, d []interface{}) error
```
CreateMapSlice creates a Etcd directory based on []interface{}.


## func CreateStruct
``` go
func CreateStruct(client *etcd.Client, dir string, s interface{}) error
```
CreateStruct creates a Etcd directory based on a struct.


## func JSON
``` go
func JSON(root *etcd.Node) ([]byte, error)
```
JSON returns an Etcd directory as JSON []byte.


## func JSONIndent
``` go
func JSONIndent(root *etcd.Node, indent string) ([]byte, error)
```
JSONIndent returns an Etcd directory as indented JSON []byte.


## func Map
``` go
func Map(root *etcd.Node) map[string]interface{}
```
Map returns a map[string]interface{} from a Etcd directory.


## func Struct
``` go
func Struct(root *etcd.Node, s interface{}) error
```
Struct returns a struct from a Etcd directory.
!!! This is not supported for nested struct yet.









- - -
