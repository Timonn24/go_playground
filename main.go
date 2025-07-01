package main

import (
	"flag"
	"fmt"
	"log"
)
// based on hashicorp/memberlist (https://pkg.go.dev/gopkg.in/hashicorp/memberlist.v0)
// examples: https://github.com/octu0/example-memberlist
// usage: <exe> -name=master -port=7947 -wport=:8080
// usage: <exe> -name=slave -port=7948 -wport=8081 -join=<master__ip>:7947
func main() {
	namePtr := flag.String("name", "nodeX", "Specify a node name")
	portPtr := flag.Int("port", 7946, "Specify a port to bind to")
	webPtr := flag.String("wport", ":8080", "Specify a web port")
	joinPtr := flag.String("join", "", "Specify cluster addr")
	flag.Parse()

	fmt.Println("Node name: ", *namePtr)
	fmt.Println("Binding port: ", *portPtr)
	fmt.Println("Web console at: ", *webPtr)
	if joinPtr != nil && len(*joinPtr) > 0 {
		fmt.Println("Cluster at: ", *joinPtr)
	}

	node := &Node{Name: *namePtr, Port: *portPtr}
	store := NewKeyValueStore()
	cluster, err := NewCluster(node, store, *joinPtr)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	api := NewAPI(cluster)
	api.Run(*webPtr)
}
