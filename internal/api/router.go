package api

import (
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
)

func NewRouter(libvirtClient *libvirt.Client, ovnClient *ovn.Client) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/libvirt/vms", createVMHandler(libvirtClient))
	mux.HandleFunc("POST /api/libvirt/vms/{name}/stop", stopVMHandler(libvirtClient))
	mux.HandleFunc("DELETE /api/libvirt/vms/{name}", deleteVMHandler(libvirtClient))

	mux.HandleFunc("POST /api/libvirt/vms/{name}/nics", attachNICHandler(libvirtClient))
	mux.HandleFunc("GET /api/libvirt/vms/{name}/nics", listNICsHandler(libvirtClient))

	mux.HandleFunc("POST /api/ovn/switches", createSwitchHandler(ovnClient))
	mux.HandleFunc("GET /api/ovn/switches/{name}", getSwitchHandler(ovnClient))
	mux.HandleFunc("DELETE /api/ovn/switches/{name}", deleteSwitchHandler(ovnClient))
	return mux
}
