package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/huseyinbabal/botkube-plugins-playground/plugin"
	plugin3 "github.com/huseyinbabal/botkube-plugins-playground/plugin/executor"
	plugin2 "github.com/huseyinbabal/botkube-plugins-playground/plugin/source"
	"log"
	"math/rand"
	"time"
)

func main() {
	go handleCloudEvents()
	go initializeSourcePlugins()
	go initializeExecutorPlugins()
	time.Sleep(time.Hour * 2)
}

func handleCloudEvents() {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("Client failure, %v", err)
	}

	err = c.StartReceiver(context.Background(), func(event cloudevents.Event) {
		fmt.Printf("%s", event)
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(err)
}

func initializeSourcePlugins() {
	sourcePlugins := plugin.NewManager("./contrib/build", plugin.TypeSource, &plugin2.SourcePlugin{})
	defer sourcePlugins.Dispose()
	if err := sourcePlugins.Initialize(); err != nil {
		log.Fatalf("failed to initialize source plugins %v", err)
	}
	if err := sourcePlugins.Start(); err != nil {
		log.Fatalf("failed to start source plugins %v", err)
	}
	events := make(chan interface{}, 1)
	for _, pl := range []string{"kubernetes"} {
		adapter, err := sourcePlugins.GetAdapter(pl)
		if err != nil {
			log.Fatal(err)
		}
		adapter.(plugin2.Source).Consume(events)
	}
	for {
		select {
		case event := <-events:
			eventMap := event.(map[string]string)
			fmt.Println("Event receved. Name: %", eventMap["Name"])
		}
	}
}

func initializeExecutorPlugins() {
	executorPlugins := plugin.NewManager("./contrib/build", plugin.TypeExecutor, &plugin3.ExecutorPlugin{})
	defer executorPlugins.Dispose()
	if err := executorPlugins.Initialize(); err != nil {
		log.Fatalf("failed to initialize executor plugins %v", err)
	}
	if err := executorPlugins.Start(); err != nil {
		log.Fatalf("failed to start executor plugins %v", err)
	}

	// Assume that we have message from slack and by using prefix, we resolve which plugin to use. e.g. kubectl
	adapter, err := executorPlugins.GetAdapter("kubectl")
	if err != nil {
		log.Fatal(err)
	}
	kubectl := adapter.(plugin3.Executor)

	commands := []string{
		"kubectl get pods",
		"kubectl describe svc kubernetes -n defaul",
		"kubectl get nodes",
	}
	for {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		cmd := commands[r.Intn(len(commands))]
		res, err := kubectl.Execute(cmd)
		if err != nil {
			fmt.Printf("error executing %s. Err: %v", cmd, res)
		} else {
			fmt.Printf("The result of %s is %v", cmd, res)
		}

		time.Sleep(time.Second * 5)
	}

}
