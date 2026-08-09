package libvirt

import (
	"fmt"
	"io"
	"net/http"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

const (
	storagePoolName = "default"

	alpineImageURL = "https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/generic_alpine-3.22.0-x86_64-bios-cloudinit-metal-r0.qcow2"
	baseImageName  = "alpine-base.qcow2"
)

type Client struct {
	conn *libvirt.Connect
}

func NewClient() (*Client, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	_, err := c.conn.Close()
	return err
}

func (c *Client) CreateVM(domain *libvirtxml.Domain) (*libvirt.Domain, error) {

	if domain.Devices == nil {
		domain.Devices = &libvirtxml.DomainDeviceList{}
	}

	pool, err := c.conn.LookupStoragePoolByName(storagePoolName)
	if err != nil {
		return nil, fmt.Errorf("lookup storage pool: %w", err)
	}
	defer pool.Free()

	baseVolume, err := ensureBaseVolume(c.conn, pool)
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

	xml, err := domain.Marshal()
	if err != nil {
		return nil, err
	}

	vm, err := c.conn.DomainDefineXML(xml)
	if err != nil {
		return nil, err
	}

	if err := vm.Create(); err != nil {
		return nil, err
	}

	return vm, nil
}

// ensureBaseVolume ensures that the Alpine base image exists as a libvirt
// storage volume.
func ensureBaseVolume(conn *libvirt.Connect, pool *libvirt.StoragePool) (*libvirt.StorageVol, error) {
	volume, err := pool.LookupStorageVolByName(baseImageName)
	if err == nil {
		return volume, nil
	}

	// The base volume does not exist yet.
	// Download the Alpine image and import it into the libvirt storage pool.
	resp, err := http.Get(alpineImageURL)
	if err != nil {
		return nil, fmt.Errorf("download alpine image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"download Alpine image: %s",
			resp.Status,
		)
	}

	volumeXML := &libvirtxml.StorageVolume{
		Name: baseImageName,
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: "qcow2",
			},
		},
	}

	xml, err := volumeXML.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal volume XML: %w", err)
	}

	volume, err = pool.StorageVolCreateXML(xml, 0)
	if err != nil {
		return nil, fmt.Errorf("create base volume: %w", err)
	}

	stream, err := conn.NewStream(0)
	if err != nil {
		volume.Delete(0)
		volume.Free()

		return nil, fmt.Errorf("create upload stream: %w", err)
	}

	if err := volume.Upload(
		stream,
		0,
		uint64(resp.ContentLength),
		0,
	); err != nil {
		stream.Abort()
		stream.Free()

		volume.Delete(0)
		volume.Free()

		return nil, fmt.Errorf("start upload base volume: %w", err)
	}

	err = stream.SendAll(func(_ *libvirt.Stream, size int) ([]byte, error) {
		buf := make([]byte, size)

		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		return buf[:n], nil
	})
	if err != nil {
		stream.Abort()
		stream.Free()

		volume.Delete(0)
		volume.Free()

		return nil, fmt.Errorf("upload base volume: %w", err)
	}

	if err := stream.Finish(); err != nil {
		stream.Free()

		volume.Delete(0)
		volume.Free()

		return nil, fmt.Errorf("finish upload base volume: %w", err)
	}

	if err := stream.Free(); err != nil {
		volume.Delete(0)
		volume.Free()

		return nil, fmt.Errorf("free upload stream: %w", err)
	}

	return volume, nil
}
