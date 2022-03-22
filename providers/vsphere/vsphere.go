package vsphere

import (
	"context"
	"fmt"
	"nb-sync/config"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

func vmwareConnect(ctx context.Context) (*govmomi.Client, error) {
	// LOAD CONFIG
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("CONFIG NOT LOADED: %s", err)
		return nil, err
	}

	// PARSE URL
	vURL := "https://" + cfg.Providers.VSphere.Username + ":" + cfg.Providers.VSphere.Password + "@" + cfg.Providers.VSphere.Host + "/sdk"
	u, err := url.Parse(vURL)
	if err != nil {
		fmt.Printf("URL PARSE ERROR: %s", u)
		return nil, err
	}

	// CONNECT TO VSPHERE
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Logging in error: %s\n", err.Error())
		return nil, err
	}

	// LOGIN SUCCESS
	fmt.Println("Log in successful")
	return c, nil
}

func VMwareSync() []string {
	fmt.Println("I'm VMware provider")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connect, err := vmwareConnect(ctx)
	if err != nil {
		fmt.Printf("CONNECT ERROR: %s\n", err)
		return nil
	}
	// FINDER
	finder := find.NewFinder(connect.Client, true)

	// GET DATACENTERS
	dataCenters, err := getDataCenters(ctx, finder)
	if err != nil {
		return nil
	}

	return dataCenters
}

func getDataCenters(ctx context.Context, finder *find.Finder) ([]string, error) {
	dcs, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		return nil, err
	}

	var data []string
	for _, i := range dcs {
		data = append(data, i.Name())
	}
	return data, nil
}

/*
func GetDatacenterss(ctx context.Context, connect *govmomi.Client) {
	finder := find.NewFinder(connect.Client, true)

	dcs, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		fmt.Printf("FIND DATACENTERS ERROR: %s\n", err)
	}

	var data []string
	for _, v := range dcs {
		data = append(data, v.InventoryPath)
	}
	fmt.Println(data)
	// df, errs = finder.FolderList(ctx, "/")
	// if errs != nil {
	// 	fmt.Printf("FIND DATACENTERS ERROR: %s\n", errs)
	// 	return
	// }
	// fmt.Println(df)

	// GET DATACENTER BY NAME
	dc, err := finder.Datacenter(ctx, "BI_DC")
	if err != nil {
		fmt.Printf("FIND DATACENTERS ERROR: %s\n", err)
		return
	}
	fmt.Println(dc)
	finder.SetDatacenter(dc)

	ccr, err := finder.ClusterComputeResourceList(ctx, "DAC")
	if err != nil {
		fmt.Printf("FIND DATACENTERS ERROR: %s\n", err)
		return
	}

	fmt.Println(ccr[0].ComputeResource.Hosts(ctx))

	// finder.VirtualMachineList(ctx)
	// fmt.Println(dc[0])
}

// func VSphereConnect() {

// 	f := find.NewFinder(c.Client, true)

// 	// Find one and only datacenter
// 	dc, err := f.Datacenter(ctx, "/BI_DC")
// 	fmt.Println("DC: %s", dc)

// }

// func VMwareGetDatacenter() {

// }
*/
