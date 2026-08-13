package core

import (
	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovs"
)

type Core struct {
	libvirt *libvirt.Client
	ovn     *ovn.Client
	ovs     *ovs.Client
}

func NewCore(libvirtClient *libvirt.Client, ovnClient *ovn.Client, ovsClient *ovs.Client) *Core {
	if libvirtClient == nil || ovnClient == nil || ovsClient == nil {
		panic("libvirtClient, ovnClient, and ovsClient must not be nil")
	}
	return &Core{
		libvirt: libvirtClient,
		ovn:     ovnClient,
		ovs:     ovsClient,
	}
}