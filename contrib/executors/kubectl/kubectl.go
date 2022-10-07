package main

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-plugin"
	plugin2 "github.com/huseyinbabal/botkube-plugins-playground/plugin"
	plugin4 "github.com/huseyinbabal/botkube-plugins-playground/plugin/executor"
	"os/exec"
	"strings"
)

type KubectlExecutor struct{}

func (KubectlExecutor) Execute(command string) (string, error) {
	return run(command)
}

func run(command string) (string, error) {
	commandParts := strings.Split(command, " ")
	cmd := exec.Command(commandParts[0], commandParts[1:]...)

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("cmd: %s, err: %s", command, errOut.String())
	}
	return out.String(), nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		Plugins: map[string]plugin.Plugin{
			"kubectl": &plugin4.ExecutorPlugin{Impl: &KubectlExecutor{}},
		},
		HandshakeConfig: plugin2.Handshake,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
