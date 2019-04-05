package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)
var (
	consulAddress = "127.0.0.1:8500"
	client *api.Client
	nodeIDs []string
)

// consul watch example
func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	var (
		err    error
		plan   *watch.Plan
		ch     chan int
	)
	ch = make(chan int, 1)

	params := make(map[string]interface{})
	params["type"] = "service"
	params["service"] = "test"
	plan, err = watch.Parse(params)
	if err != nil {
		panic(err)
	}

	// This block here shows where the type being sent is the vendored []*api.ServiceEntry, not the accessible one
	plan.Handler = func(index uint64, result interface{}) {
		if entries, ok := result.([]*api.ServiceEntry); ok {
			// We expect to see this trigger
			for _, e := range entries {
				fmt.Printf("serviceEntry: %#v\n", e.Service)
			}
			ch <- 1
		} else {
			// We always see this trigger
			fmt.Printf("watch data was not of the expected type. Got %v\n", reflect.TypeOf(result).String())
		}
	}

	go func() {
		if err = plan.Run(consulAddress); err != nil {
			panic(err)
		}
	}()
	go http.ListenAndServe(":8080", nil)
	for i := 0; i < 2; i++ {
		id := register()
		fmt.Printf("registered %s\n", id)
		nodeIDs = append(nodeIDs, id)
	}
	for {
		<-ch
		fmt.Printf("got a change in the service entries\n")
	}
}

func cleanup() {
	for _, id := range nodeIDs {
		fmt.Printf("deregistering %s\n", id)
		if err := client.Agent().ServiceDeregister(id); err != nil {
			fmt.Println(err.Error())
		}
	}

	// Give some time for consul to deregister
	time.Sleep(2*time.Second)
}

func register() string {
	var (
		err    error
	)
	client, err = api.NewClient(&api.Config{Address:consulAddress})
	if err != nil {
		panic(err)
	}
	id := randomString()
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:   id,
		Name: "test",
		Tags: []string{"SERVER"},
		Port: 8080,
		Check: &api.AgentServiceCheck{
			HTTP: "",
		},
	})
	if err != nil {
		panic(err)
	}

	return id
}

func randomString() string {
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}
