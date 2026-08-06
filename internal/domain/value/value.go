package value

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/netip"
	"strings"
	"time"
)

var ErrInvalidValue = errors.New("invalid value")

const ulidAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

type IPv4 struct {
	addr netip.Addr
}

func NewIPv4(raw string) (IPv4, error) {
	addr, err := netip.ParseAddr(strings.TrimSpace(raw))
	if err != nil || !addr.Is4() {
		return IPv4{}, fmt.Errorf("%w: ipv4 %q", ErrInvalidValue, raw)
	}
	return IPv4{addr: addr}, nil
}

func (v IPv4) String() string {
	if !v.addr.IsValid() {
		return ""
	}
	return v.addr.String()
}

func (v IPv4) IsZero() bool {
	return !v.addr.IsValid()
}

func (v IPv4) Addr() netip.Addr {
	return v.addr
}

type IPv6 struct {
	addr netip.Addr
}

func NewIPv6(raw string) (IPv6, error) {
	addr, err := netip.ParseAddr(strings.TrimSpace(raw))
	if err != nil || !addr.Is6() {
		return IPv6{}, fmt.Errorf("%w: ipv6 %q", ErrInvalidValue, raw)
	}
	return IPv6{addr: addr}, nil
}

func (v IPv6) String() string {
	if !v.addr.IsValid() {
		return ""
	}
	return v.addr.String()
}

func (v IPv6) IsZero() bool {
	return !v.addr.IsValid()
}

func (v IPv6) Addr() netip.Addr {
	return v.addr
}

type CIDR struct {
	prefix netip.Prefix
}

func NewCIDR(raw string) (CIDR, error) {
	prefix, err := netip.ParsePrefix(strings.TrimSpace(raw))
	if err != nil || !prefix.IsValid() {
		return CIDR{}, fmt.Errorf("%w: cidr %q", ErrInvalidValue, raw)
	}
	return CIDR{prefix: prefix}, nil
}

func (c CIDR) String() string {
	if !c.prefix.IsValid() {
		return ""
	}
	return c.prefix.String()
}

func (c CIDR) IsZero() bool {
	return !c.prefix.IsValid()
}

func (c CIDR) Prefix() netip.Prefix {
	return c.prefix
}

func (c CIDR) Contains(addr netip.Addr) bool {
	return c.prefix.IsValid() && c.prefix.Contains(addr)
}

func (c CIDR) IsIPv4() bool {
	return c.prefix.IsValid() && c.prefix.Addr().Is4()
}

func (c CIDR) IsIPv6() bool {
	return c.prefix.IsValid() && c.prefix.Addr().Is6()
}

type MACAddress string

func NewMACAddress(raw string) (MACAddress, error) {
	addr, err := net.ParseMAC(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: mac %q", ErrInvalidValue, raw)
	}
	return MACAddress(strings.ToLower(addr.String())), nil
}

func (m MACAddress) String() string {
	return string(m)
}

func (m MACAddress) IsZero() bool {
	return strings.TrimSpace(string(m)) == ""
}

func NewPrefixedULID(prefix string) string {
	timestamp := uint64(time.Now().UTC().UnixMilli())
	var entropy [10]byte
	if _, err := rand.Read(entropy[:]); err != nil {
		for i := range entropy {
			entropy[i] = byte(timestamp >> uint((i%8)*8))
		}
	}

	var raw [16]byte
	raw[0] = byte(timestamp >> 40)
	raw[1] = byte(timestamp >> 32)
	raw[2] = byte(timestamp >> 24)
	raw[3] = byte(timestamp >> 16)
	raw[4] = byte(timestamp >> 8)
	raw[5] = byte(timestamp)
	copy(raw[6:], entropy[:])

	encoded := encodeBase32(raw[:])
	if prefix == "" {
		return encoded
	}
	return prefix + "-" + encoded
}

func encodeBase32(raw []byte) string {
	value := new(big.Int).SetBytes(raw)
	base := big.NewInt(32)
	zero := big.NewInt(0)
	chars := make([]byte, 26)
	for i := 25; i >= 0; i-- {
		if value.Cmp(zero) == 0 {
			chars[i] = ulidAlphabet[0]
			continue
		}
		mod := new(big.Int)
		value.DivMod(value, base, mod)
		chars[i] = ulidAlphabet[mod.Int64()]
	}
	return string(chars)
}

func CloneStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return map[string]string{}
	}
	clone := make(map[string]string, len(src))
	for key, value := range src {
		clone[key] = value
	}
	return clone
}

func CloneByteMap(src map[string][]byte) map[string][]byte {
	if len(src) == 0 {
		return map[string][]byte{}
	}
	clone := make(map[string][]byte, len(src))
	for key, value := range src {
		copied := make([]byte, len(value))
		copy(copied, value)
		clone[key] = copied
	}
	return clone
}
