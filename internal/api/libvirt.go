package api

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
// It is registered as POST /api/libvirt/vms
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

// stopVMHandler handles VM stop requests.
//
// It is registered as POST /api/libvirt/vms/{name}/stop
func stopVMHandler(client *libvirt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/api/libvirt/vms/") : len(r.URL.Path)-len("/stop")]
		if err := client.StopVM(name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// deleteVMHandler handles VM deletion requests.
//
// It is registered as DELETE /api/libvirt/vms/{name}
func deleteVMHandler(client *libvirt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/api/libvirt/vms/"):]
		if err := client.DeleteVM(name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}