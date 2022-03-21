package providers

import (
	"fmt"
	"nb-sync/providers/vmware"
)

func ProviderInit() {
	fmt.Println("Hello")
	vmware.VSphereConnect()
}
