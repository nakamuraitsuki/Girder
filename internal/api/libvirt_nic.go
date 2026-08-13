package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"libvirt.org/go/libvirtxml"
)

// AttachNICRequest represents the NIC atacch parameters accepted by the HTTP API.
type AttachNICRequest struct {
	LogicalName string `json:"logicalName"`
}

// NICResponse represents a NIC as exporsed by the HTTP API.
type NICResponse struct {
	LogicalName string `json:"logicalName"`
	TapDevice   string `json:"tapDevice"`
	MAC         string `json:"mac"`
}

// attachNICHandler handles NIC attach requests.
//
// It is registered as POST /api/libvirt/vms/{name}/nics
func attachNICHandler(client *libvirt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		var req AttachNICRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(req.LogicalName, "ua-") {
			req.LogicalName = "ua-" + req.LogicalName
		}

		iface := &libvirtxml.DomainInterface{
			Alias: &libvirtxml.DomainAlias{
				Name: req.LogicalName,
			},
		}

		actual, err := client.AttachNIC(name, iface)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(toNICResponse(*actual))
	}
}

// listNICsHandler handles NIC list requests.
//
// It is registerd as GET /api/libvirt/vms/{name}/nics
func listNICsHandler(client *libvirt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		ifaces, err := client.ListNICs(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := make([]NICResponse, 0, len(ifaces))
		for _, iface := range ifaces {
			resp = append(resp, toNICResponse(iface))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// toNICResponse is a helper function that converts a libvirt DomainInterface to a NICResponse.
func toNICResponse(iface libvirtxml.DomainInterface) NICResponse {
	resp := NICResponse{
		TapDevice: iface.Target.Dev,
	}
	if iface.Alias != nil {
		resp.LogicalName = iface.Alias.Name
	}
	if iface.MAC != nil {
		resp.MAC = iface.MAC.Address
	}
	return resp
}
