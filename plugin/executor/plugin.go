package plugin

import (
	"context"
	"github.com/hashicorp/go-plugin"
	executor "github.com/huseyinbabal/botkube-plugins-playground/plugin/executor/proto"
	"google.golang.org/grpc"
)

type ExecutorPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	plugin.GRPCPlugin
	Impl Executor
}

func (p *ExecutorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	executor.RegisterExecutorServer(s, &ExecutorGRPCServer{
		Impl:   p.Impl,
		Broker: broker,
	})
	return nil
}

func (p *ExecutorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ExecutorGRPCClient{
		Client: executor.NewExecutorClient(c),
		Broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &ExecutorPlugin{}
