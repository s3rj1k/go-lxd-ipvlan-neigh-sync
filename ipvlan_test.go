package main

import (
	"net"
	"reflect"
	"testing"

	"github.com/lxc/lxd/shared/api"
)

func TestExtractIPVlanNeigh(t *testing.T) {
	tests := []struct {
		name       string
		linkName   string              // in
		containers []api.ContainerFull // in
		want       *IPVlanNeigh        // out
	}{
		{
			name:     "1 IPv4",
			linkName: "vlan10",
			containers: []api.ContainerFull{
				{
					Container: api.Container{
						ContainerPut: api.ContainerPut{
							Devices: map[string]map[string]string{
								"eth0": {
									"name":         "eth0",
									"nictype":      "ipvlan",
									"parent":       "vlan10",
									"type":         "nic",
									"ipv4.address": "192.0.2.101",
									"mtu":          "1500",
								},
							},
						},
					},
					State: &api.ContainerState{
						Pid: 3851793,
					},
				},
			},
			want: &IPVlanNeigh{
				LinkName: "vlan10",
				IP: map[string]net.IP{
					"192.0.2.101": net.ParseIP("192.0.2.101"),
				},
			},
		},
		{
			name:     "2 IPv4",
			linkName: "vlan10",
			containers: []api.ContainerFull{
				{
					Container: api.Container{
						ContainerPut: api.ContainerPut{
							Devices: map[string]map[string]string{
								"eth0": {
									"ipv4.address": "192.0.2.101,192.0.2.102",
									"mtu":          "1500",
									"name":         "eth0",
									"nictype":      "ipvlan",
									"parent":       "vlan10",
									"type":         "nic",
								},
							},
						},
					},
					State: &api.ContainerState{
						Pid: 3851793,
					},
				},
			},
			want: &IPVlanNeigh{
				LinkName: "vlan10",
				IP: map[string]net.IP{
					"192.0.2.101": net.ParseIP("192.0.2.101"),
					"192.0.2.102": net.ParseIP("192.0.2.102"),
				},
			},
		},
		{
			name:     "IPv4 with spaces",
			linkName: "vlan10",
			containers: []api.ContainerFull{
				{
					Container: api.Container{
						ContainerPut: api.ContainerPut{
							Devices: map[string]map[string]string{
								"eth0": {
									"ipv4.address": "192.0.2.101, 192.0.2.102, 192.0.2.103,192.0.2.104,",
									"name":         "eth0",
									"nictype":      "ipvlan",
									"parent":       "vlan10",
									"type":         "nic",
								},
							},
						},
					},
					State: &api.ContainerState{
						Pid: 1000,
					},
				},
			},
			want: &IPVlanNeigh{
				LinkName: "vlan10",
				IP: map[string]net.IP{
					"192.0.2.101": net.ParseIP("192.0.2.101"),
					"192.0.2.102": net.ParseIP("192.0.2.102"),
					"192.0.2.103": net.ParseIP("192.0.2.103"),
					"192.0.2.104": net.ParseIP("192.0.2.104"),
				},
			},
		},
		{
			name:     "IPv6 & IPv4",
			linkName: "vlan10",
			containers: []api.ContainerFull{
				{
					Container: api.Container{
						ContainerPut: api.ContainerPut{
							Devices: map[string]map[string]string{
								"eth0": {
									"name":         "eth0",
									"nictype":      "ipvlan",
									"type":         "nic",
									"ipv6.address": "2001:db8::101",
									"mtu":          "1500",
									"parent":       "vlan10",
									"ipv4.address": "192.0.2.101",
								},
							},
						},
					},
					State: &api.ContainerState{
						Pid: 434980,
					},
				},
			},
			want: &IPVlanNeigh{
				LinkName: "vlan10",
				IP: map[string]net.IP{
					"2001:db8::101": net.ParseIP("2001:db8::101"),
					"192.0.2.101":   net.ParseIP("192.0.2.101"),
				},
			},
		},
		{
			name:     "2 Containers",
			linkName: "vlan10",
			containers: []api.ContainerFull{
				{
					Container: api.Container{
						ContainerPut: api.ContainerPut{
							Devices: map[string]map[string]string{
								"eth0": {
									"ipv4.address": "192.0.2.101",
									"mtu":          "1500",
									"name":         "eth0",
									"nictype":      "ipvlan",
									"parent":       "vlan10",
									"type":         "nic",
								},
							},
						},
					},
					State: &api.ContainerState{
						Pid: 1000,
					},
				},
				{
					Container: api.Container{
						ContainerPut: api.ContainerPut{
							Devices: map[string]map[string]string{
								"eth0": {
									"ipv4.address": "192.0.2.102",
									"mtu":          "1500",
									"name":         "eth0",
									"nictype":      "ipvlan",
									"parent":       "vlan10",
									"type":         "nic",
								},
							},
						},
					},
					State: &api.ContainerState{
						Pid: 2000,
					},
				},
			},
			want: &IPVlanNeigh{
				LinkName: "vlan10",
				IP: map[string]net.IP{
					"192.0.2.101": net.ParseIP("192.0.2.101"),
					"192.0.2.102": net.ParseIP("192.0.2.102"),
				},
			},
		},
	}

	for _, testCase := range tests {
		tt := testCase

		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractIPVlanNeigh(tt.linkName, tt.containers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractIPVlanNeigh() = %v, want %v", got, tt.want)
			}
		})
	}
}
