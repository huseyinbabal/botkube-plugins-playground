package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/huseyinbabal/botkube-plugins-playground/plugin"
	botkubeexecutorplugin "github.com/huseyinbabal/botkube-plugins/api/executor"
	botkubesourceplugin "github.com/huseyinbabal/botkube-plugins/api/source"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	go handleCloudEvents()
	go initializePlugins()
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

func initializePlugins() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to find user home folder. %v", err)
	}
	plugins := plugin.NewManager(filepath.Join(homeDir, ".botkube-plugins-cache"))
	err = plugins.Initialize([]string{"kubernetes", "kubectl"})
	if err != nil {
		log.Fatalf("failedt to initialize plugins %v", err)
	}
	defer plugins.Dispose()
	if err := plugins.Start(); err != nil {
		log.Fatalf("failed to start plugins %v", err)
	}

	go initializeSourcePlugins(plugins)
	go initializeExecutorPlugins(plugins)
}

func initializeSourcePlugins(manager *plugin.Manager) {
	events := make(chan interface{}, 1)
	// Let say that we have only kubernetes plugin enabled
	for _, pl := range []string{"kubernetes"} {
		adapter, err := manager.GetAdapter(pl)
		if err != nil {
			log.Fatal(err)
		}
		adapter.(botkubesourceplugin.Source).Consume(events)
	}
	for {
		select {
		case event := <-events:
			eventMap := event.(map[string]string)
			fmt.Println("Event received. Name: %", eventMap["Name"])
		}
	}
}

func initializeExecutorPlugins(manager *plugin.Manager) {

	// Assume that we have message from slack and by using prefix, we resolve which plugin to use. e.g. kubectl
	adapter, err := manager.GetAdapter("kubectl")
	if err != nil {
		log.Fatal(err)
	}
	kubectl := adapter.(botkubeexecutorplugin.Executor)

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
