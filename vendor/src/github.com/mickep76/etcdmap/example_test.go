package etcdmap_test

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	Name      string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Group struct {
	Name  string `json:"groupname"`
	Users []User `json:"users"`
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
				FirstName: "John",
				LastName:  "Doe",
			},
			User{
				Name:      "lnemoy",
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
	err := etcdmap.CreateStruct(client, "/example", g)
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
	//{"groupname":"staff","users":{"0":{"first_name":"John","last_name":"Doe","username":"jdoe"},"1":{"first_name":"Leonard","last_name":"Nimoy","username":"lnemoy"}}}
}
