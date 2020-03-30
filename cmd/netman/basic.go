package main

import (
	"context"

	"github.com/pkg/errors"
	config "go.ligato.io/vpp-agent/v3/client"
	agent "go.ligato.io/vpp-agent/v3/cmd/agentctl/client"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
)

// TestLoopbackGRPC creation
func TestLoopbackGRPC(client config.ConfigClient) (err error) {
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
	changeReq := client.ChangeRequest()
	changeReq.Update(loop)

	return changeReq.Send(context.TODO())
}

// TestLoopbackETCD creation
func TestLoopbackETCD(client *agent.Client) (err error) {
	// Create a LOOPBACK interface using the ETCD client
	loop := &vpp_interfaces.Interface{
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
	if kvdb, err = client.KVDBClient(); err != nil {
		panic(errors.Wrap(err, "syncClient.KVDBClient()"))
	}
	var key string
	if key, err = kvdb.CompleteFullKey(models.Key(loop)); err != nil {
		panic(errors.Wrap(err, "kvdb.CompleteFullKey()"))
	}
	return kvdb.ProtoBroker().Put(key, loop)
}
