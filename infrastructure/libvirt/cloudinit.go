package libvirt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func createCloudInitSeed(domainName string) (string, error) {
	dir := "/var/lib/libvirt/images"

	userDataPath := filepath.Join(dir, domainName+"-user-data")
	metaDataPath := filepath.Join(dir, domainName+"-meta-data")
	seedPath := filepath.Join(dir, domainName+"-seed.iso")

	userData := `#cloud-config
password: root
chpasswd:
  expire: false
ssh_pwauth: true
`

	metaData := fmt.Sprintf(`instance-id: %s
local-hostname: %s
`, domainName, domainName)

	if err := os.WriteFile(userDataPath, []byte(userData), 0644); err != nil {
		return "", fmt.Errorf("write user-data: %w", err)
	}
	defer os.Remove(userDataPath)

	if err := os.WriteFile(metaDataPath, []byte(metaData), 0644); err != nil {
		return "", fmt.Errorf("write meta-data: %w", err)
	}
	defer os.Remove(metaDataPath)

	cmd := exec.Command(
		"cloud-localds",
		seedPath,
		userDataPath,
		metaDataPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf(
			"create cloud-init seed: %w: %s",
			err,
			output,
		)
	}

	return seedPath, nil
}