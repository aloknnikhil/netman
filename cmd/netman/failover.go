package main

import (
	"context"

	config "go.ligato.io/vpp-agent/v3/client"
	vpp_interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	vpp_nat "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/nat"
)

// TestFailover of WAN interfaces
func TestFailover(client config.ConfigClient) (err error) {

	// Start change transaction
	// ATOMIC
	changeReq := client.ChangeRequest()

	// Define WAN 0
	wan0 := &vpp_interfaces.Interface{
		Name:          "wan0",
		Type:          vpp_interfaces.Interface_DPDK,
		SetDhcpClient: true,
		Vrf:           1,
	}
	changeReq.Update(wan0)

	// Define LAN 1
	lan1 := &vpp_interfaces.Interface{
		Name:          "lan1",
		Type:          vpp_interfaces.Interface_DPDK,
		SetDhcpClient: true,
		Vrf:           0,
	}
	changeReq.Update(lan1)

	// Setup NAT interface in/out
	natWan0 := &vpp_nat.Nat44Interface{
		Name:       "wan0",
		NatOutside: false,
		NatInside:  false,
	}
	changeReq.Delete(natWan0)

	natLan1 := &vpp_nat.Nat44Interface{
		Name:       "lan1",
		NatOutside: true,
		NatInside:  false,
	}
	changeReq.Update(natLan1)

	// Setup NAT Indentity Mapping
	dnat := &vpp_nat.DNat44{
		Label: "Egress",
		IdMappings: []*vpp_nat.DNat44_IdentityMapping{
			&vpp_nat.DNat44_IdentityMapping{
				VrfId:     0,
				Interface: "lan1",
			},
		},
	}
	changeReq.Update(dnat)

	// Triggering a failover should NOT be persisted to ETCD. Use GRPC to trigger a one-off sequence of commands
	return changeReq.Send(context.TODO())
}
