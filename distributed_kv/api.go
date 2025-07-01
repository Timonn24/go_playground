package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type API struct {
	router  *mux.Router
	cluster *Cluster
}

func NewAPI(cluster *Cluster) *API {
	api := &API{router: mux.NewRouter(), cluster: cluster}
	api.setupRoutes()
	return api
}

func (api *API) setupRoutes() {
	api.router.HandleFunc("/set/{key}/{value}", api.setHandler).Methods("POST")
	api.router.HandleFunc("/get/{key}", api.getHandler).Methods("GET")
}

func (api *API) setHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("setHandler called!")
	parts := strings.Split(r.URL.Path, "/")
	api.cluster.store.Set(parts[2], parts[3])

	tempMap := make(map[string]any)
	tempMap[parts[2]] = parts[3]
	data, err := json.Marshal(&tempMap)
	if err == nil {
		for _, node := range api.cluster.Memberlist.Members() {
			if node.Name == api.cluster.LocalNode.Name {
				continue
			}
			api.cluster.Memberlist.SendReliable(node, data)
		}
	}
}

func (api *API) getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getHandler called!")
	key_pos := strings.LastIndex(r.URL.Path, "/")
	fmt.Println(api.cluster.store.Get(r.URL.Path[key_pos+1:]))

}

func (api *API) Run(addr string) {
	go func() {
		run := true
		for run {
			select {
			case data := <-api.cluster.LocalNode.msgChan:
				err := api.cluster.store.FromBytes(data)
				if err != nil {
					continue
				}

				log.Printf("Received message. Data size = %d", len(data))
			}
		}
	}()
	http.ListenAndServe(addr, api.router)
}
