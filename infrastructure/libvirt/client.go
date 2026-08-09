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
		vm.Undefine()
		vm.Free()
		return nil, err
	}

	return vm, nil
}

// ensureBaseVolume ensures that the Alpine base image exists as a libvirt
// storage volume.
func (c *Client) ensureBaseVolume(pool *libvirt.StoragePool) (*libvirt.StorageVol, error) {
	volume, err := pool.LookupStorageVolByName(baseImageName)
	if err == nil {
		return volume, nil
	}

	resp, err := http.Get(alpineImageURL)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	if resp.ContentLength <= 0 {
		return nil, fmt.Errorf("unknown content length, cannot pre-declare capacity")
	}

	volXML := &libvirtxml.StorageVolume{
		Name:     baseImageName,
		Capacity: &libvirtxml.StorageVolumeSize{Value: uint64(resp.ContentLength), Unit: "bytes"},
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{Type: "qcow2"},
		},
	}
	volXMLText, err := volXML.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal volume XML: %w", err)
	}

	vol, err := pool.StorageVolCreateXML(volXMLText, 0)
	if err != nil {
		return nil, fmt.Errorf("create volume: %w", err)
	}

	if err := c.UploadVolume(vol, resp.Body, uint64(resp.ContentLength)); err != nil {
		vol.Free()
		return nil, err
	}

	return vol, nil
}

// UploadVolume writes the contents of r into vol via a libvirt stream.
// It normalizes libvirt's Stream.Send-based API to the standard io.Reader
// interface so callers never need to know about virStream semantics.
func (c *Client) UploadVolume(vol *libvirt.StorageVol, r io.Reader, length uint64) error {
	stream, err := c.conn.NewStream(0)
	if err != nil {
		return fmt.Errorf("new stream: %w", err)
	}
	defer stream.Free()

	if err := vol.Upload(stream, 0, length, 0); err != nil {
		return fmt.Errorf("upload: %w", err)
	}

	if _, err := io.Copy(streamWriter{stream}, r); err != nil {
		stream.Abort()
		return fmt.Errorf("stream copy: %w", err)
	}

	return stream.Finish()
}

// streamWriter adapts *libvirt.Stream to io.Writer.
// This is unexported and confined to this file: it exists only because
// Stream.Send has the io.Writer signature but not the io.Writer name.
type streamWriter struct {
	stream *libvirt.Stream
}

func (w streamWriter) Write(p []byte) (int, error) {
	return w.stream.Send(p)
}