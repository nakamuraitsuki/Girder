package http

import (
	"encoding/json"
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"libvirt.org/go/libvirtxml"
)

// CreateVMRequest represents the general VM parameters accepted by the HTTP API.
type CreateVMRequest struct {
	Name   string `json:"name"`
	Memory uint `json:"memory"`
	VCPUs  uint   `json:"vcpus"`
}

// createVMHandler handles VM creation requests.
//
// The handler maps the general VM representation used by the HTTP API to the
// libvirt-specific Domain representation before passing it to the libvirt client.
func createVMHandler(client *libvirt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateVMRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		domain := libvirtxml.Domain{
			Type: "kvm",
			Name: req.Name,
			Memory: &libvirtxml.DomainMemory{
				Value: req.Memory,
				Unit:  "KiB",
			},
			VCPU: &libvirtxml.DomainVCPU{
				Value: req.VCPUs,
			},
		}

		if _, err := client.CreateVM(&domain); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
