package api

import (
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
)

func NewRouter(libvirtClient *libvirt.Client) http.Handler {
	mux := http.NewServeMux()
	
	mux.HandleFunc("POST /api/libvirt/vms", createVMHandler(libvirtClient))
	mux.HandleFunc("POST /api/libvirt/vms/{name}/stop", stopVMHandler(libvirtClient))
	mux.HandleFunc("DELETE /api/libvirt/vms/{name}", deleteVMHandler(libvirtClient))
	return mux
}