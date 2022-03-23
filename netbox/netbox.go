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
		fmt.Printf("NETBOX NOT CONNECT: %s", connect)
		return nil
	}
	fmt.Println("NETBOX CONNECTED")

	return connect
}

// Основная функция синхронизации данных из VSPHERE
func SyncData(dcs []string) {
	// SYNC TAG
	syncTagExist := SyncTagFind()
	if !syncTagExist {
		SyncTagCreate()
	}

	ClusterGroupSync(dcs)
}

// Функция синхронизации ClusterGroup(DataCenters)
// Принимает в себя массив dcs: новых CG из VSphere и eDcs: массив существующих CG
// TODO: Скорее всего эту функцию нужно будет вынести в другой пакет: providers
func ClusterGroupSync(dcs []string) {
	// Исключить из синхронизации CG из конфига (Exclude)
	exArr := excludeFilter(dcs, cfg.Settings.DataCenters.Exclude)

	// TODO: Включить в синхронизацию CG из конфига (Include)
	// FIXME: inArr := include(exArr, dcs)

	existCgs := ClusterGroupList()
	var cg []string
	for _, o := range existCgs {
		cg = append(cg, *o.Name)
	}

	// Сравнить 2 массива (exArr и eDcs) на выходе должны получить значения которые не имеются во 2м массиве
	addDcs, deleteDcs := diffData(exArr, cg)

	// Create ClusterGroup
	for _, s := range addDcs {
		ClusterGroupCreate(s, s)
	}
	for _, s1 := range deleteDcs {
		for _, s2 := range existCgs {
			if s1 == *s2.Name {
				ClusterGroupDelete(s2.ID)
			}
		}
	}
	fmt.Println("SYNCED")
}

// Создать новый ClusterGroup в Netbox
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
		log.Fatalf("CREATE CLUSTERGROUP REQUEST ERROR: %s", err)
		return nil
	}

	return clusterGroup
}

// Удалить существующий ClusterGroup в Netbox
func ClusterGroupDelete(cgID int64) *virtualization.VirtualizationClusterGroupsDeleteNoContent {
	params := &virtualization.VirtualizationClusterGroupsDeleteParams{
		ID:      cgID,
		Context: context.Background(),
	}
	clusterGroup, err := connect.Virtualization.VirtualizationClusterGroupsDelete(params, nil)
	if err != nil {
		log.Fatalf("DELETE CLUSTERGROUP REQUEST ERROR: %s", err)
		return nil
	}

	return clusterGroup
}

// Получить существующие ClusterGroup в Netbox
func ClusterGroupList() []*models.ClusterGroup {
	cgs, err := connect.Virtualization.VirtualizationClusterGroupsList(nil, nil)
	if err != nil {
		log.Fatalf("CREATE ClusterGroupList REQUEST ERROR: %s", err)
		return nil
	}

	return cgs.Payload.Results
}

// Создать тег системы синхронизации в Netbox
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
		log.Fatalf("CREATE SYNCTAG REQUEST ERROR: %s", err)
		return
	}
	syncTagID = syncTag.Payload.ID
}

// Проверка существования тега системы синхронизации в Netbox
func SyncTagFind() bool {
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

/* HELPERS */
// TODO: Вынести все эти функции в отдельный пакет: tools

// Исключает массив из массива
// TODO: Вернуть вторым параметром значения которые были исключены
func excludeFilter(arr []string, arr2 []string) []string {
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

// Включить в массив значения которые описаны в конфиге
func includeFilter(arr []string, arr2 []string) []string {

	return nil
}

// Получить различия между двумя слайсами
func diffData(arr1 []string, arr2 []string) ([]string, []string) {
	var dataAdd []string
	var dataDelete []string

	for i := 0; i < 2; i++ {
		for _, s1 := range arr1 {
			found := false
			for _, s2 := range arr2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found && i == 0 {
				dataAdd = append(dataAdd, s1)
			} else if !found && i == 1 {
				dataDelete = append(dataDelete, s1)
			}
		}
		if i == 0 {
			arr1, arr2 = arr2, arr1
		}
	}
	return dataAdd, dataDelete
}
