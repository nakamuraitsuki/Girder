package ovn

import (
	"crypto/rand"
	"fmt"
)

// generateMAC generates a random locally-administered unicast MAC address.
//
// The first octet has bit 1 (locally administered) set and bit 0
// (multicast) cleared, following the same convention libvirt/OVN
// tooling uses for auto-assigned addresses (e.g. 52:54:00:.. style,
// but here we use 02:.. to avoid collisions with libvirt's default OUI).
func generateMAC() (string, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate random bytes for MAC address: %w", err)
	}

	buf[0] = (buf[0] | 0x02) & 0xfe // Set locally administered bit, clear multicast bit

	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5],
	), nil
}
