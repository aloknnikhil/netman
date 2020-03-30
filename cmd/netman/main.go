package main

import (
	"context"

	"go.ligato.io/vpp-agent/v3/pkg/models"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"

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

	// Create a LOOPBACK interface using the GRPC client
	loop := &vpp_interfaces.Interface{
		Name:        "loopMeGRPC",
		Type:        vpp_interfaces.Interface_SOFTWARE_LOOPBACK,
		Enabled:     true,
		PhysAddress: "de:ad:be:ef:ba:ad",
		IpAddresses: []string{
			"1.2.3.4/32",
			"2.3.4.5/32",
		},
		Mtu: 1500,
	}

	// Start change transaction
	// ATOMIC
	// Add all necessary steps here before "sending" the transaction
	changeReq := configClient.ChangeRequest()
	changeReq.Update(loop)

	if err = changeReq.Send(context.TODO()); err != nil {
		panic(errors.Wrap(err, "changeReq.Send()"))
	}

	// Create a LOOPBACK interface using the ETCD client
	loop = &vpp_interfaces.Interface{
		Name:        "loopMeETCD",
		Type:        vpp_interfaces.Interface_SOFTWARE_LOOPBACK,
		Enabled:     true,
		PhysAddress: "de:ad:be:ef:99:99",
		IpAddresses: []string{
			"11.12.13.14/32",
			"12.13.14.15/32",
		},
		Mtu: 1500,
	}

	var kvdb agent.KVDBAPIClient
	if kvdb, err = syncClient.KVDBClient(); err != nil {
		panic(errors.Wrap(err, "syncClient.KVDBClient()"))
	}
	var key string
	if key, err = kvdb.CompleteFullKey(models.Key(loop)); err != nil {
		panic(errors.Wrap(err, "kvdb.CompleteFullKey()"))
	}
	kvdb.ProtoBroker().Put(key, loop)
}
