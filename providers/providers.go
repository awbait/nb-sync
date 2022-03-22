package providers

import (
	"nb-sync/netbox"
	"nb-sync/providers/vsphere"
)

func ProviderInit() {
	netbox.SyncData(vsphere.VMwareSync())
}
