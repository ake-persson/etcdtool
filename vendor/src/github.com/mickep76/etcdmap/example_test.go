package etcdmap_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcd "github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
)

// Env variables.
type Env struct {
	Peers string
}

// User structure.
type User struct {
	Name      string `json:"username" etcd:"name"`
	Age       int    `json:"age" etcd:"age"`
	Male      bool   `json:"male" etcd:"male"`
	FirstName string `json:"first_name" etcd:"first_name"`
	LastName  string `json:"last_name" etcd:"last_name"`
}

// Group structure.
type Group struct {
	Name  string `json:"groupname" etcd:"name"`
	Users []User `json:"users" etcd:"users"`
}

// getEnv variables.
func getEnv() Env {
	env := Env{}
	env.Peers = "http://127.0.0.1:4001,http://127.0.0.1:2379"

	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		switch a[0] {
		case "ETCD_PEERS":
			env.Peers = a[1]
		}
	}

	return env
}

// ExampleNestedStruct creates a Etcd directory using a nested Go struct and then gets the directory as JSON.
func Example_nestedStruct() {
	// Get env variables.
	env := getEnv()

	// Options.
	peers := flag.String("peers", env.Peers, "Comma separated list of etcd nodes, can be set with env. variable ETCD_PEERS")
	flag.Parse()

	// Define nested structure.
	g := Group{
		Name: "staff",
		Users: []User{
			User{
				Name:      "jdoe",
				Age:       25,
				Male:      true,
				FirstName: "John",
				LastName:  "Doe",
			},
			User{
				Name:      "lnemoy",
				Age:       62,
				Male:      true,
				FirstName: "Leonard",
				LastName:  "Nimoy",
			},
		},
	}

	// Connect to etcd.
	cfg := etcd.Config{
		Endpoints:               strings.Split(*peers, ","),
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	client, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := etcd.NewKeysAPI(client)

	// Create directory structure based on struct.
	err2 := etcdmap.Create(kapi, "/example", reflect.ValueOf(g))
	if err2 != nil {
		log.Fatal(err2.Error())
	}

	// Get directory structure from Etcd.
	res, err3 := kapi.Get(context.TODO(), "/example", &etcd.GetOptions{Recursive: true})
	if err3 != nil {
		log.Fatal(err3.Error())
	}

	j, err4 := etcdmap.JSON(res.Node)
	if err4 != nil {
		log.Fatal(err4.Error())
	}

	fmt.Println(string(j))

	// Get directory structure from Etcd.
	/*
		res, err5 := kapi.Get(context.TODO(), "/example/users/0", &etcd.GetOptions{Recursive: true})
		if err5 != nil {
			log.Fatal(err5.Error())
		}

		s := User{}
		err6 := etcdmap.Struct(res.Node, reflect.ValueOf(&s))
		if err6 != nil {
			log.Fatal(err6.Error())
		}

		fmt.Println(s)
	*/

	// Output:
	//{"name":"staff","users":{"0":{"age":"25","first_name":"John","last_name":"Doe","male":"true","name":"jdoe"},"1":{"age":"62","first_name":"Leonard","last_name":"Nimoy","male":"true","name":"lnemoy"}}}
}
