package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func main() {
	vURL := flag.String("url", "", "The URL of a vCenter server")
	flag.Parse()

	u, err := url.Parse(*vURL)
	if err != nil {
		fmt.Printf("Error parsing url %s\n", vURL)
		return
	}

	fmt.Println("%s", u.User)

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
	// Make future calls local to this datacenter
	f.SetDatacenter(dc)

	// Find virtual machines in datacenter
	vms, err := f.VirtualMachineList(ctx, "*")

	pc := property.DefaultCollector(c.Client)

	// Convert datastores into list of references
	var refs []types.ManagedObjectReference
	for _, vm := range vms {
		refs = append(refs, vm.Reference())
	}

	// Retrieve name property for all vms
	var vmt []mo.VirtualMachine
	err = pc.Retrieve(ctx, refs, []string{"name"}, &vmt)

	// Print name per virtual machine
	tw := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Println("Virtual machines found:", len(vmt))
	// sort.Sort(ByName(vmt))
	for _, vm := range vmt {
		fmt.Fprintf(tw, "%s\n", vm.Name)
	}
	tw.Flush()

	c.Logout(ctx)
}
