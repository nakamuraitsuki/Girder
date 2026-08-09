package libvirt

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

	baseVolume, err := ensureBaseVolume(pool)
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
func ensureBaseVolume(pool *libvirt.StoragePool) (*libvirt.StorageVol, error) {
	volume, err := pool.LookupStorageVolByName(baseImageName)
	if err == nil {
		return volume, nil
	}

	// default pool must be a directory-based storage pool.
	poolXML, err := pool.GetXMLDesc(0)
	if err != nil {
		return nil, fmt.Errorf("get storage pool XML: %w", err)
	}

	var poolDef libvirtxml.StoragePool
	if err := poolDef.Unmarshal(poolXML); err != nil {
		return nil, fmt.Errorf("unmarshal storage pool XML: %w", err)
	}

	if poolDef.Target == nil || poolDef.Target.Path == "" {
		return nil, fmt.Errorf("storage pool %s has no target path", storagePoolName)
	}

	imagePath := filepath.Join(
		poolDef.Target.Path,
		baseImageName,
	)

	// Download the base image if it dowe not exist yet.
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		if err := downloadFile(alpineImageURL, imagePath); err != nil {
			return nil, fmt.Errorf("download base image: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("check base image existence: %w", err)
	}

	// Tell libvirt to rescan the directory pool.
	if err := pool.Refresh(0); err != nil {
		return nil, fmt.Errorf("refresh storage pool: %w", err)
	}

	volume, err = pool.LookupStorageVolByName(baseImageName)
	if err != nil {
		return nil, fmt.Errorf("lookup base volume after refresh: %w", err)
	}

	return volume, nil
}

// downloadFile downloads a file from the given URL and saves it to the specified
//
// This function is utility function so it is tail of the file.
func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http get: unexpected status code %d", resp.StatusCode)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("copy response body to file: %w", err)
	}

	return nil
}