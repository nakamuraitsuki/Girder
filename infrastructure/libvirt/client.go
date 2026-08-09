package libvirt

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

const (
	imageDir = "/var/lib/libvirt/images"

	alpineImageURL =
		"https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/generic_alpine-3.22.0-x86_64-bios-cloudinit-metal-r0.qcow2"
	baseImageName = "alpine-base.qcow2"
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
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, err
	}

	baseImage := filepath.Join(imageDir, baseImageName)
	diskPath := filepath.Join(
		imageDir,
		fmt.Sprintf("%s.qcow2", domain.Name),
	)

	if err := ensureBaseImage(baseImage); err != nil {
		return nil, err
	}

	if err := createDisk(baseImage, diskPath); err != nil {
		return nil, err
	}

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
					File: diskPath,
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

func ensureBaseImage(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	resp, err := http.Get(alpineImageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to download Alpine image: %s",
			resp.Status,
		)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func createDisk(baseImage, diskPath string) error {
	if _, err := os.Stat(diskPath); err == nil {
		return fmt.Errorf("disk already exists: %s", diskPath)
	} else if !os.IsNotExist(err) {
		return err
	}

	cmd := exec.Command(
		"qemu-img",
		"create",
		"-f",
		"qcow2",
		"-F",
		"qcow2",
		"-b",
		baseImage,
		diskPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"qemu-img failed: %w: %s",
			err,
			string(output),
		)
	}

	return nil
}