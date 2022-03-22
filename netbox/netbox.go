package netbox

import (
	"context"
	"fmt"
	"log"
	"nb-sync/config"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func netboxConnect() *client.NetBoxAPI {
	// LOAD CONFIG
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("CONFIG NOT LOADED: %s", err)
		return nil
	}

	connect := netbox.NewNetboxWithAPIKey(fmt.Sprintf("%s:%d", cfg.Netbox.Host, cfg.Netbox.Port), cfg.Netbox.Token)
	if connect == nil {
		fmt.Printf("NETBOX NOT CONNECT: %s", connect)
		return nil
	}
	fmt.Println("NETBOX CONNECTED")

	return connect
}

func SyncData(dcs []string) {
	connect := netboxConnect()

	for _, s := range dcs {
		CreateClusterGroup(connect, s, s)
	}
}

func CreateClusterGroup(connect *client.NetBoxAPI, name string, slug string) *virtualization.VirtualizationClusterGroupsCreateCreated {
	params := &virtualization.VirtualizationClusterGroupsCreateParams{
		Data: &models.ClusterGroup{
			Name: &name,
			Slug: &slug,
		},
		Context: context.Background(),
	}
	clusterGroup, err := connect.Virtualization.VirtualizationClusterGroupsCreate(params, nil)
	if err != nil {
		log.Fatalf("CREATE CLUSTERGROUP REQUEST ERROR: %s", err)
		return nil
	}

	return clusterGroup
}

// TEST
// gl, err := connect.Virtualization.VirtualizationClusterGroupsList(nil, nil)
// if err != nil {
// 	fmt.Printf("REQUEST ERR: %s", err)
// 	return nil
// }

// fmt.Println(*gl.Payload.Results[0].Name)
// #TEST
