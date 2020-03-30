package main

import (
	"github.com/pkg/errors"
	config "go.ligato.io/vpp-agent/v3/client"
	agent "go.ligato.io/vpp-agent/v3/cmd/agentctl/client"
)

const (
	// DefaultAgentHost defines default host address for agent.
	DefaultAgentHost = "127.0.0.1"
	// DefaultPortGRPC defines default port for GRPC connection.
	DefaultPortGRPC = 9111
	// DefaultETCDEndpoint defines default endpoint for ETCD connection.
	DefaultETCDEndpoint = "0.0.0.0:2379"
	// DefaultServiceLabel of the VPP instance to control
	DefaultServiceLabel = "meter_vpp"
)

func main() {
	var err error
	var syncClient *agent.Client
	if syncClient, err = agent.NewClientWithOpts(
		agent.WithServiceLabel(DefaultServiceLabel),
		agent.WithHost(DefaultAgentHost),
		agent.WithGrpcPort(DefaultPortGRPC),
		agent.WithEtcdEndpoints(
			[]string{
				DefaultETCDEndpoint,
			},
		),
	); err != nil {
		panic(errors.Wrap(err, "agent.NewClientWithOpts()"))
	}

	var configClient config.ConfigClient
	if configClient, err = syncClient.ConfigClient(); err != nil {
		panic(errors.Wrap(err, "agent.ConfigClient()"))
	}

	// TestLoopbackGRPC
	if err = TestLoopbackGRPC(configClient); err != nil {
		panic(errors.Wrap(err, "TestLoopbackGRPC()"))
	}

	// TestLoopbackETCD
	if err = TestLoopbackETCD(syncClient); err != nil {
		panic(errors.Wrap(err, "TestLoopbackETCD()"))
	}

	// TestFailover
	if err = TestFailover(configClient); err != nil {
		panic(errors.Wrap(err, "TestFailover()"))
	}
}
