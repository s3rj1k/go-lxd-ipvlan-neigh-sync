package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/vishvananda/netlink"
)

/*
	https://github.com/lxc/lxc/issues/3445

	ip neigh add proxy {IPv4} dev {IfaceName}
	ip neigh list proxy
	ip neigh del proxy {IPv4} dev {IfaceName}

	ip -6 neigh add proxy {IPv6} dev {IfaceName}
	ip -6 neigh list proxy
	ip -6 neigh del proxy {IPv6} dev {IfaceName}
*/

// IPVlanNeigh desribes ProxyARP(NDP) table entries.
type IPVlanNeigh struct {
	LinkName string

	IP map[string]net.IP
}

// InitIPVlanNeigh returns initialized ProxyARP(NDP) object.
func InitIPVlanNeigh() (neigh *IPVlanNeigh) {
	out := new(IPVlanNeigh)

	out.IP = make(map[string]net.IP)

	return out
}

// GetLinkIndex returns network interface index by name.
func (neigh *IPVlanNeigh) GetLinkIndex() (int, error) {
	link, err := netlink.LinkByName(neigh.LinkName)
	if err != nil {
		return 0, fmt.Errorf("netlink data error: %w", err)
	}

	return link.Attrs().Index, nil
}

// GetIPVlanNeigh returns ProxyARP(NDP) table entries.
func GetIPVlanNeigh(linkName string) (*IPVlanNeigh, error) {
	// create ProxyARP(NDP) table object
	out := InitIPVlanNeigh()

	// set network interface name
	out.LinkName = strings.TrimSpace(linkName)

	// set connection properties
	connArgs := &lxd.ConnectionArgs{
		// custom HTTP Client (used as base for the connection)
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	// connect to LXD over the Unix socket ("/var/lib/lxd/unix.socket")
	conn, err := lxd.ConnectLXDUnix("", connArgs)
	if err != nil {
		return nil, fmt.Errorf("lxd connection error: %w", err)
	}

	// check LXD extension presence
	if !conn.HasExtension("container_full") {
		return nil, fmt.Errorf("lxd connection error: no container_full API extension")
	}

	// contains all containers information
	var containers []api.ContainerFull

	// lxc query -X GET --wait /1.0/containers?recursion=2
	containers, err = conn.GetContainersFull()
	if err != nil {
		return nil, fmt.Errorf("lxd data error: %w", err)
	}

	// process all containers information
	for _, ct := range containers {
		// skip container that does not have any init system PID
		if ct.State.Pid == 0 {
			continue
		}

		// process all container devices
		for _, dev := range ct.Devices {
			// skip none network device
			if val := dev["type"]; val != "nic" {
				continue
			}

			// skip none IPVLAN network type
			if val := dev["nictype"]; val != "ipvlan" {
				continue
			}

			var (
				parentLink, v4Address, v6Address string

				ok bool
			)

			// get parent network interface name
			parentLink, ok = dev["parent"]
			if !ok {
				continue
			}

			// skip unspecified interface
			if !strings.EqualFold(
				out.LinkName,
				strings.TrimSpace(parentLink),
			) {
				continue
			}

			// get IPv4 addresses
			v4Address, ok = dev["ipv4.address"]
			if !ok {
				continue
			}

			// append IPv4 to output table
			for _, el := range strings.Fields(v4Address) {
				if ip := net.ParseIP(el); ip != nil {
					out.IP[ip.String()] = ip
				}
			}

			// get IPv6 addresses
			v6Address, ok = dev["ipv6.address"]
			if !ok {
				continue
			}

			// append IPv6 to output table
			for _, el := range strings.Fields(v6Address) {
				if ip := net.ParseIP(el); ip != nil {
					out.IP[ip.String()] = ip
				}
			}
		}
	}

	return out, nil
}
