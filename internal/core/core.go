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
