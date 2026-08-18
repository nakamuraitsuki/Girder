package libvirt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func (c *Client) createCloudInitVolume(
	pool *libvirt.StoragePool,
	domainName string,
) (*libvirt.StorageVol, error) {
	volumeName := fmt.Sprintf("%s-seed.iso", domainName)

	// Reuse an existing seed volume if present.
	volume, err := pool.LookupStorageVolByName(volumeName)
	if err == nil {
		return volume, nil
	}

	tmpDir, err := os.MkdirTemp("", "girder-cloud-init-*")
	if err != nil {
		return nil, fmt.Errorf("create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	userDataPath := filepath.Join(tmpDir, "user-data")
	metaDataPath := filepath.Join(tmpDir, "meta-data")
	seedPath := filepath.Join(tmpDir, "seed.iso")

	userData := `#cloud-config
password: root
chpasswd:
  expire: false
ssh_pwauth: true
`

	metaData := fmt.Sprintf(`instance-id: %s
local-hostname: %s
`, domainName, domainName)

	if err := os.WriteFile(userDataPath, []byte(userData), 0600); err != nil {
		return nil, fmt.Errorf("write user-data: %w", err)
	}

	if err := os.WriteFile(metaDataPath, []byte(metaData), 0600); err != nil {
		return nil, fmt.Errorf("write meta-data: %w", err)
	}

	cmd := exec.Command(
		"cloud-localds",
		seedPath,
		userDataPath,
		metaDataPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf(
			"create cloud-init image: %w: %s",
			err,
			output,
		)
	}

	info, err := os.Stat(seedPath)
	if err != nil {
		return nil, fmt.Errorf("stat cloud-init image: %w", err)
	}

	volXML := &libvirtxml.StorageVolume{
		Name: volumeName,
		Capacity: &libvirtxml.StorageVolumeSize{
			Value: uint64(info.Size()),
			Unit:  "bytes",
		},
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: "raw",
			},
		},
	}

	volXMLText, err := volXML.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal cloud-init volume XML: %w", err)
	}

	volume, err = pool.StorageVolCreateXML(volXMLText, 0)
	if err != nil {
		return nil, fmt.Errorf("create cloud-init volume: %w", err)
	}

	file, err := os.Open(seedPath)
	if err != nil {
		volume.Delete(0)
		volume.Free()
		return nil, fmt.Errorf("open cloud-init image: %w", err)
	}
	defer file.Close()

	if err := c.UploadVolume(volume, file, uint64(info.Size())); err != nil {
		volume.Delete(0)
		volume.Free()
		return nil, fmt.Errorf("upload cloud-init volume: %w", err)
	}

	return volume, nil
}