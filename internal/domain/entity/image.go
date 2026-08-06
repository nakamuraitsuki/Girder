package entity

import (
	"fmt"
	"strings"
)

type ImageState string

const (
	ImageStateActive   ImageState = "active"
	ImageStateDisabled ImageState = "disabled"
)

type Image struct {
	ID          ImageID
	Name        string
	Description string
	OS          string
	Architecture string
	Version     string
	SizeBytes   int64
	Checksum    string
	State       ImageState
	Tags        map[string]string
}

func NewImage(name, osName, architecture, version string, sizeBytes int64) (Image, error) {
	if strings.TrimSpace(name) == "" {
		return Image{}, fmt.Errorf("%w: image name", ErrEmptyName)
	}
	if sizeBytes < 0 {
		return Image{}, fmt.Errorf("%w: image size", ErrInvalidResource)
	}
	return Image{
		ID:          NewImageID(),
		Name:        strings.TrimSpace(name),
		OS:          strings.TrimSpace(osName),
		Architecture: strings.TrimSpace(architecture),
		Version:     strings.TrimSpace(version),
		SizeBytes:   sizeBytes,
		State:       ImageStateActive,
		Tags:        map[string]string{},
	}, nil
}

func (i *Image) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: image name", ErrEmptyName)
	}
	i.Name = strings.TrimSpace(name)
	return nil
}

func (i *Image) SetDescription(description string) {
	i.Description = strings.TrimSpace(description)
}

func (i *Image) SetVersion(version string) {
	i.Version = strings.TrimSpace(version)
}

func (i *Image) SetChecksum(checksum string) {
	i.Checksum = strings.TrimSpace(checksum)
}

func (i *Image) Enable() {
	i.State = ImageStateActive
}

func (i *Image) Disable() {
	i.State = ImageStateDisabled
}

func (i *Image) SetTag(key, val string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("%w: tag key", ErrEmptyName)
	}
	if i.Tags == nil {
		i.Tags = map[string]string{}
	}
	i.Tags[key] = strings.TrimSpace(val)
	return nil
}

func (i *Image) RemoveTag(key string) {
	if i.Tags == nil {
		return
	}
	delete(i.Tags, strings.TrimSpace(key))
}