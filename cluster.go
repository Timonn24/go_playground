package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/memberlist"
)

type Node struct {
	Name    string
	Port    int
	Addr    string
	msgChan chan []byte
}

type Cluster struct {
	*memberlist.Memberlist
	LocalNode *Node
	store     *KeyValueStore
}

type Delegate struct {
	msgChan chan []byte
}

func (d *Delegate) NotifyMsg(msg []byte) {
	fmt.Println("Got a message: ", msg)
	d.msgChan <- msg
}

func (d *Delegate) NodeMeta(limit int) []byte {
	return nil
}

func (d *Delegate) LocalState(join bool) []byte {
	return nil
}

func (d *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d *Delegate) MergeRemoteState(buf []byte, join bool) {
}

func NewCluster(localNode *Node, store *KeyValueStore, cluster_addr string) (*Cluster, error) {

	msgChan := make(chan []byte)
	d := new(Delegate)
	d.msgChan = msgChan

	config := memberlist.DefaultLocalConfig()
	config.Name = localNode.Name
	//config.BindAddr = localNode.Addr
	config.BindPort = localNode.Port
	config.Delegate = d

	localNode.msgChan = d.msgChan

	list, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}

	local := list.LocalNode()
	list.Join([]string{
		fmt.Sprintf("%s:%d", local.Addr.To4().String(), local.Port),
	})

	if len(cluster_addr) > 0 {
		log.Printf("cluster join to %s", cluster_addr)
		if _, err := list.Join([]string{cluster_addr}); err != nil {
			log.Fatal(err)
		}
	}

	return &Cluster{Memberlist: list, LocalNode: localNode, store: store}, nil
}
