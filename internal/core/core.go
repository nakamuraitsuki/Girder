package core

import (
	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
)

type Core struct {
	libvirt *libvirt.Client
	ovn     *ovn.Client
	// TODO: ovs Client implementation
}
