package etcdmap_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/mickep76/etcdmap"
)

func getEnv() []string {
	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		if a[0] == "ETCD_CONN" {
			return []string{a[1]}
		}
	}

	return []string{}
}

type User struct {
	Name      string `json:"username" etcd:"id"`
	Age       int    `json:"age" etcd:"age"`
	Male      bool   `json:"male" etcd:"male"`
	FirstName string `json:"first_name" etcd:"first_name"`
	LastName  string `json:"last_name" etcd:"last_name"`
}

type Group struct {
	Name  string `json:"groupname" etcd:"id"`
	Users []User `json:"users" etcd:"users"`
}

// ExampleNestedStruct creates a Etcd directory using a nested Go struct and then gets the directory as JSON.
func Example_nestedStruct() {
	verbose := flag.Bool("verbose", false, "Verbose")
	node := flag.String("node", "", "Etcd node")
	port := flag.String("port", "", "Etcd port")
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

	// Connect to Etcd.
	conn := getEnv()
	if node == nil && port == nil {
		conn = []string{fmt.Sprintf("http://%v:%v", *node, *port)}
	}

	if *verbose {
		log.Printf("Connecting to: %s", conn)
	}
	client := etcd.NewClient(conn)

	// Create directory structure based on struct.
	err := etcdmap.Create(client, "/example", reflect.ValueOf(g))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get directory structure from Etcd.
	res, err := client.Get("/example", true, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	j, err2 := etcdmap.JSON(res.Node)
	if err2 != nil {
		log.Fatal(err2.Error())
	}

	fmt.Println(string(j))

	// Output:
	//{"id":"staff","users":{"0":{"age":"25","first_name":"John","id":"jdoe","last_name":"Doe","male":"true"},"1":{"age":"62","first_name":"Leonard","id":"lnemoy","last_name":"Nimoy","male":"true"}}}
}
