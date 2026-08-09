package libvirt

import (
	"fmt"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// CreateVM creates a new virtual machine using the provided libvirt domain.
//
// 1. It ensures that the base volume is present in the default storage pool.
// 2. It creates a new storage volume for the VM based on the base volume.
// 3. It adds default disk and console devices to the domain.
// 4. It defines and starts the VM using the libvirt client.
func (c *Client) CreateVM(domain *libvirtxml.Domain) (*libvirt.Domain, error) {

	// ensure the default storage pool exists and retrieve it
	pool, err := c.conn.LookupStoragePoolByName(storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("lookup storage pool: %w", err)
	}
	defer pool.Free()

	baseVolume, err := c.ensureBaseVolume(pool)
	if err != nil {
		return nil, fmt.Errorf("prepare base bolume: %w", err)
	}
	defer baseVolume.Free()

	volumeName := fmt.Sprintf("%s.qcow2", domain.Name)

	volumeXML := &libvirtxml.StorageVolume{
		Name: volumeName,
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: "qcow2",
			},
		},
	}

	volumeXMLText, err := volumeXML.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal volume XML: %w", err)
	}

	volume, err := pool.StorageVolCreateXMLFrom(volumeXMLText, baseVolume, 0)
	if err != nil {
		return nil, fmt.Errorf("create volume from base: %w", err)
	}
	defer volume.Free()

	volumePath, err := volume.GetPath()
	if err != nil {
		return nil, fmt.Errorf("get volume path: %w", err)
	}

	// Adds default disk and console devices to the domain if they are not already present.
	if domain.Devices == nil {
		domain.Devices = &libvirtxml.DomainDeviceList{}
	}

	domain.Devices.Disks = append(
		domain.Devices.Disks,
		libvirtxml.DomainDisk{
			Device: "disk",
			Driver: &libvirtxml.DomainDiskDriver{
				Name: "qemu",
				Type: "qcow2",
			},
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: volumePath,
				},
			},
			Target: &libvirtxml.DomainDiskTarget{
				Dev: "vda",
				Bus: "virtio",
			},
		},
	)

	domain.OS = &libvirtxml.DomainOS{
		Type: &libvirtxml.DomainOSType{
			Arch:    "x86_64",
			Machine: "q35",
			Type:    "hvm",
		},
		BootDevices: []libvirtxml.DomainBootDevice{
			{
				Dev: "hd",
			},
		},
	}

	port := uint(0)

	domain.Devices.Consoles = append(
		domain.Devices.Consoles,
		libvirtxml.DomainConsole{
			Source: &libvirtxml.DomainChardevSource{
				Pty: &libvirtxml.DomainChardevSourcePty{},
			},
			Target: &libvirtxml.DomainConsoleTarget{
				Type: "serial",
				Port: &port,
			},
		},
	)

	// Define and start the VM using prepared domain XML. 
	xml, err := domain.Marshal()
	if err != nil {
		return nil, err
	}

	vm, err := c.conn.DomainDefineXML(xml)
	if err != nil {
		return nil, err
	}

	if err := vm.Create(); err != nil {
		vm.Undefine()
		vm.Free()
		return nil, err
	}

	return vm, nil
}
