package main

func worker(linkName string) {
	var (
		neigh   *IPVlanNeigh
		ifIndex int

		err error
	)

	// get neighbour table
	neigh, err = GetIPVlanNeigh(linkName)
	if err != nil {
		Fatal.Printf("%s\n", err)
	}

	// get interface index
	ifIndex, err = neigh.GetLinkIndex()
	if err != nil {
		Fatal.Printf("%s\n", err)
	}

	// set IP neighbour table entries
	for _, ip := range neigh.IP {
		if err := NeighSet(ifIndex, ip); err != nil {
			Error.Printf("%s\n", err)
		}
	}

	// get IP neighbour table entries
	table, err := NeighProxyList(ifIndex)
	if err != nil {
		Fatal.Printf("%s\n", err)
	}

	// process neighbour table entries
	for _, el := range table {
		if filterNeighEntry(el, ifIndex) {
			continue
		}

		// remove invalid entries
		if _, ok := neigh.IP[el.IP.String()]; !ok {
			if err := NeighDel(ifIndex, el.IP); err != nil {
				Error.Printf("%s\n", err)
			}
		}
	}
}
