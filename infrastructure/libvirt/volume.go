package libvirt

import (
	"fmt"
	"io"
	"net/http"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

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
