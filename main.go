package main

import (
	"os"

	"github.com/vishvananda/netlink"
)

func main() {
	// disble all flags for journald pretty-print
	Fatal.SetFlags(0)

	if err := CheckIfRunUnderRoot(); err != nil {
		Fatal.Fatalf("%s\n", err)
	}

	if len(os.Args) > 1 {
		Fatal.Fatalf("wrong number of arguments\n")
	}

	// netlink interface update chanel
	ch := make(chan netlink.LinkUpdate)
	if err := netlink.LinkSubscribe(ch, nil); err != nil {
		Fatal.Fatal(err)
	}

	// watch chanel events
	for update := range ch {
		attrs := update.Link.Attrs()

		if attrs.OperState == netlink.OperUp {
			Info.Printf("NIC '%s' changed its status to UP, syncing LXD proxy neighbours table\n", attrs.Name)

			go worker(attrs.Name)
		}
	}
}
