package plugin

import (
	"context"
	"github.com/hashicorp/go-plugin"
	executor "github.com/huseyinbabal/botkube-plugins-playground/plugin/executor/proto"
)

type ExecutorGRPCServer struct {
	Impl   Executor
	Broker *plugin.GRPCBroker
	executor.UnimplementedExecutorServer
}

func (p *ExecutorGRPCServer) Execute(ctx context.Context, request *executor.ExecuteRequest) (*executor.ExecuteResponse, error) {
	result, err := p.Impl.Execute(request.Command)
	if err != nil {
		return nil, err
	}
	return &executor.ExecuteResponse{Data: result}, nil
}

type ExecutorGRPCClient struct {
	Broker *plugin.GRPCBroker
	Client executor.ExecutorClient
}

func (p *ExecutorGRPCClient) Execute(command string) (string, error) {
	res, err := p.Client.Execute(context.Background(), &executor.ExecuteRequest{Command: command})
	if err != nil {
		return "", err
	}
	return res.Data, nil
}
