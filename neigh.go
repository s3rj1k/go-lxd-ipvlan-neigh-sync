package main

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// NeighSet sets single proxy neighbour entry.
func NeighSet(ifIndex int, ip net.IP) error {
	var family = netlink.FAMILY_V4

	if IsIPv4(ip) {
		family = netlink.FAMILY_V4
	}

	if IsIPv6(ip) {
		family = netlink.FAMILY_V6
	}

	if err := netlink.NeighSet(
		&netlink.Neigh{
			LinkIndex: ifIndex,
			Family:    family,
			State:     netlink.NUD_NONE,
			Type:      netlink.NDA_UNSPEC,
			Flags:     netlink.NTF_PROXY,
			IP:        ip,
		},
	); err != nil {
		return fmt.Errorf("netlink action error, ifIndex=%d, ip=%s: %w", ifIndex, ip, err)
	}

	return nil
}

// NeighDel deletes single proxy neighbour entry.
func NeighDel(ifIndex int, ip net.IP) error {
	var family = netlink.FAMILY_V4

	if IsIPv4(ip) {
		family = netlink.FAMILY_V4
	}

	if IsIPv6(ip) {
		family = netlink.FAMILY_V6
	}

	if err := netlink.NeighDel(
		&netlink.Neigh{
			LinkIndex: ifIndex,
			Family:    family,
			State:     netlink.NUD_NONE,
			Type:      netlink.NDA_UNSPEC,
			Flags:     netlink.NTF_PROXY,
			IP:        ip,
		},
	); err != nil {
		return fmt.Errorf("netlink action error, ifIndex=%d, ip=%s: %w", ifIndex, ip, err)
	}

	return nil
}

// NeighProxyList returns proxy neighbour table.
func NeighProxyList(ifIndex int) ([]netlink.Neigh, error) {
	// get IPv4 neighbour table entries
	neigh4, err := netlink.NeighProxyList(ifIndex, netlink.FAMILY_V4)
	if err != nil {
		return nil, fmt.Errorf("netlink data error: %w", err)
	}

	// get IPv6 neighbour table entries
	neigh6, err := netlink.NeighProxyList(ifIndex, netlink.FAMILY_V6)
	if err != nil {
		return nil, fmt.Errorf("netlink data error: %w", err)
	}

	return append(neigh4, neigh6...), nil
}

func filterNeighEntry(neigh netlink.Neigh, ifIndex int) bool {
	if neigh.LinkIndex != ifIndex { // check interface
		return true
	}

	if neigh.State != netlink.NUD_NONE { // check state
		return true
	}

	if neigh.Type != netlink.NDA_UNSPEC { // check type
		return true
	}

	if neigh.Flags != netlink.NTF_PROXY { // check flags
		return true
	}

	return false
}
