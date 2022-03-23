package netbox

import (
	"context"
	"fmt"
	"log"
	"nb-sync/config"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/extras"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

var syncTagSlug string = "nb_sync"
var syncTagID int64

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

	// SYNC TAG
	syncTagExist := SyncTagFind(connect)
	if !syncTagExist {
		SyncTagCreate(connect)
	}

	// eDcs := ClusterGroupList(connect)

	// ClusterGroupSync(dcs, eDcs)
	// for _, s := range dcs {
	// 	ClusterGroupCreate(connect, s, s)
	// }
}

func ClusterGroupSync(dcs []string, eDcs []*models.ClusterGroup) {
	// При вызове функции получаем массив с нашими датацентрами
	// Далее, получаем список датацентров из Netbox'a

	//
	// LOAD CONFIG
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("CONFIG NOT LOADED: %s", err)
		return
	}
	// END DELETE

	// Исключаем из синхронизации
	farr := exlude(dcs, cfg.Settings.DataCenters.Exclude)
	fmt.Println(farr)

	// TODO: Включить в синхронизацию если такое есть в наличии из конфига
	// TODO: Получить список ClusterGroup из netbox'a и
}

func ClusterGroupCreate(connect *client.NetBoxAPI, name string, slug string) *virtualization.VirtualizationClusterGroupsCreateCreated {
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

func ClusterGroupList(connect *client.NetBoxAPI) []*models.ClusterGroup {
	cgs, err := connect.Virtualization.VirtualizationClusterGroupsList(nil, nil)
	if err != nil {
		log.Fatalf("CREATE ClusterGroupList REQUEST ERROR: %s", err)
		return nil
	}

	return cgs.Payload.Results
}

func SyncTagCreate(connect *client.NetBoxAPI) {
	name := "NB-Sync"
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
		log.Fatalf("CREATE SYNCTAG REQUEST ERROR: %s", err)
		return
	}
	syncTagID = syncTag.Payload.ID
}

func SyncTagFind(connect *client.NetBoxAPI) bool {
	tags, err := connect.Extras.ExtrasTagsList(nil, nil)
	if err != nil {
		log.Fatalf("LIST SYNCTAGFIND REQUEST ERROR: %s", err)
		return false
	}

	for _, tag := range tags.Payload.Results {
		if *tag.Slug == syncTagSlug {
			// TAG EXIST
			syncTagID = tag.ID
			return true
		}
	}

	// TAG NOT EXIST
	return false
}

// func SyncTagAddToItem(connect *client.NetBoxAPI) {
// 	connect.Virtualization.Tag
// }

/* HELPERS */

func exlude(arr []string, arr2 []string) []string {
	for i := 0; i < len(arr); i++ {
		el := arr[i]
		for _, rem := range arr2 {
			if el == rem {
				arr = append(arr[:i], arr[i+1:]...)
				i-- // Important: decrease index
				break
			}
		}
	}
	return arr
}

/*
	arr: Массив с исключенными значениями
	arr2: Массив с включенными значениями
*/
func include(arr []string, arr2 []string) []string {
	return nil
}

/*
TEST
gl, err := connect.Virtualization.VirtualizationClusterGroupsList(nil, nil)
if err != nil {
	fmt.Printf("REQUEST ERR: %s", err)
	return nil
}

fmt.Println(*gl.Payload.Results[0].Name)
#TEST


fmt.Println(cgs.Payload.Results)
for _, v := range cgs.Payload.Results {
	fmt.Println(*&v.ID)
}
*/
