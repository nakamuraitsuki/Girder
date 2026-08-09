package http

import (
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
)

func NewRouter(libvirtClient *libvirt.Client) http.Handler {
	mux := http.NewServeMux()
	
	mux.HandleFunc("POST /api/libvirt/vms", createVMHandler(libvirtClient))
	return mux
}