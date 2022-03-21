package vmware

import (
	"context"
	"fmt"
	"nb-sync/config"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

func VSphereConnect() {
	fmt.Println("I'm VMware provider")

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("ERR1: %s", err)
	}
	fmt.Println("User: %s", cfg.Providers.VSphere.Username)

	vURL := "https://" + cfg.Providers.VSphere.Username + ":" + cfg.Providers.VSphere.Password + "@" + cfg.Providers.VSphere.Host + "/sdk"
	u, err := url.Parse(vURL)
	if err != nil {
		fmt.Printf("Error parsing url %s\n", u)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Logging in error: %s\n", err.Error())
		return
	}

	fmt.Println("Log in successful")

	f := find.NewFinder(c.Client, true)

	// Find one and only datacenter
	dc, err := f.Datacenter(ctx, "/BI_DC")
	fmt.Println("DC: %s", dc)

}

func VMwareGetDatacenter() {
	
}