package netbox

import (
	"context"
	"fmt"
	"nb-sync/config"
	"nb-sync/tools/log"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/extras"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

var syncTagName string = "NB-Sync"
var syncTagSlug string = "nb_sync"
var syncTagID int64

var connect *client.NetBoxAPI
var cfg config.Config

func init() {
	cfg = *config.GetConfig()
	connect = netboxConnect()
}

// Функция подключения к NetBox
func netboxConnect() *client.NetBoxAPI {
	connect := netbox.NewNetboxWithAPIKey(fmt.Sprintf("%s:%d", cfg.Netbox.Host, cfg.Netbox.Port), cfg.Netbox.Token)
	if connect == nil {
		log.Error.Println("NETBOX: Not connect", connect)
		return nil
	}
	log.Info.Println("NETBOX: Connected")

	return connect
}

// Создать новый ClusterGroup
func ClusterGroupCreate(name string, slug string) *virtualization.VirtualizationClusterGroupsCreateCreated {
	params := &virtualization.VirtualizationClusterGroupsCreateParams{
		Data: &models.ClusterGroup{
			Name: &name,
			Slug: &slug,
		},
		Context: context.Background(),
	}
	clusterGroup, err := connect.Virtualization.VirtualizationClusterGroupsCreate(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: ClusterGroupCreate: ", err)
		return nil
	}

	return clusterGroup
}

// Удалить существующий ClusterGroup
func ClusterGroupDelete(cgID int64) *virtualization.VirtualizationClusterGroupsDeleteNoContent {
	params := &virtualization.VirtualizationClusterGroupsDeleteParams{
		ID:      cgID,
		Context: context.Background(),
	}
	clusterGroup, err := connect.Virtualization.VirtualizationClusterGroupsDelete(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: ClusterGroupDelete: ", err)
		return nil
	}

	return clusterGroup
}

// Получить существующие ClusterGroup
func ClusterGroupList() []*models.ClusterGroup {
	cgs, err := connect.Virtualization.VirtualizationClusterGroupsList(nil, nil)
	if err != nil {
		log.Error.Println("NETBOX: ClusterGroupList: ", err)
		return nil
	}

	return cgs.Payload.Results
}

// Создать новый ClusterType
func ClusterTypeCreate(name string, slug string) *virtualization.VirtualizationClusterTypesCreateCreated {
	params := &virtualization.VirtualizationClusterTypesCreateParams{
		Data: &models.ClusterType{
			Name: &name,
			Slug: &slug,
		},
		Context: context.Background(),
	}
	clusterType, err := connect.Virtualization.VirtualizationClusterTypesCreate(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: ClusterTypeCreate: ", err)
		return nil
	}

	return clusterType
}

// Проверить существует ли ClusterType
func ClusterTypeCheck(name string) bool {
	params := &virtualization.VirtualizationClusterTypesListParams{
		Name:    &name,
		Context: context.Background(),
	}

	clusterType, err := connect.Virtualization.VirtualizationClusterTypesList(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: ClusterTypeCheck: ", err)
		return false
	}
	if *clusterType.Payload.Count == 0 {
		return false
	}
	return true
}

// Создать новый Cluster
func ClusterCreate(name string, clusterTypeID int64, clusterGroupID int64) *virtualization.VirtualizationClustersCreateCreated {
	params := &virtualization.VirtualizationClustersCreateParams{
		Data: &models.WritableCluster{
			Name:  &name,
			Type:  &clusterTypeID,
			Group: &clusterGroupID,
			Tags:  []*models.NestedTag{{Name: &syncTagName, Slug: &syncTagSlug}},
		},
		Context: context.Background(),
	}
	cluster, err := connect.Virtualization.VirtualizationClustersCreate(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: ClusterCreate: ", err)
		return nil
	}

	return cluster
}

// Получить существующие Cluster
func ClusterList() []*models.Cluster {
	params := &virtualization.VirtualizationClustersListParams{
		Tag:     &syncTagSlug,
		Context: context.Background(),
	}
	clusters, err := connect.Virtualization.VirtualizationClustersList(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: ClusterList: ", err)
		return nil
	}

	return clusters.Payload.Results
}

// Создать новую VM
func VmCreate(name string, clusterID int64) *virtualization.VirtualizationVirtualMachinesCreateCreated {
	params := &virtualization.VirtualizationVirtualMachinesCreateParams{
		Data: &models.WritableVirtualMachineWithConfigContext{
			Name:    &name,
			Cluster: &clusterID,
			Tags:    []*models.NestedTag{{Name: &syncTagName, Slug: &syncTagSlug}},
		},
		Context: context.Background(),
	}

	vm, err := connect.Virtualization.VirtualizationVirtualMachinesCreate(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: VmCreate: ", err)
		return nil
	}

	return vm
}

// Получить существующие VM
func VmList() []*models.VirtualMachineWithConfigContext {
	params := &virtualization.VirtualizationVirtualMachinesListParams{
		Tag:     &syncTagSlug,
		Context: context.Background(),
	}
	vms, err := connect.Virtualization.VirtualizationVirtualMachinesList(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: VmList: ", err)
		return nil
	}

	return vms.Payload.Results
}

// Создать тег системы синхронизации
func SyncTagCreate() {
	name := syncTagName
	slug := syncTagSlug
	description := "DO NOT DELETE: NetBox sync system tag."
	params := &extras.ExtrasTagsCreateParams{
		Data: &models.Tag{
			Name:        &name,
			Slug:        &slug,
			Description: description,
		},
		Context: context.Background(),
	}
	syncTag, err := connect.Extras.ExtrasTagsCreate(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: SyncTagCreate: ", err)
		return
	}
	syncTagID = syncTag.Payload.ID
}

// Проверка существования тега системы синхронизации
func SyncTagCheck() bool {
	params := &extras.ExtrasTagsListParams{
		Slug: &syncTagSlug,
	}
	tags, err := connect.Extras.ExtrasTagsList(params, nil)
	if err != nil {
		log.Error.Println("NETBOX: SyncTagCheck: ", err)
		return false
	}

	if *tags.Payload.Count == 0 {
		return false
	}

	return true
}
