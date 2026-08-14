package api

import (
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovs"
	"github.com/nakamuraitsuki/Girder/internal/core"
)

func NewRouter(
	libvirtClient *libvirt.Client,
	ovnClient *ovn.Client,
	ovsClient *ovs.Client,
	core *core.Core,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/libvirt/vms", createVMHandler(libvirtClient))
	mux.HandleFunc("POST /api/libvirt/vms/{name}/stop", stopVMHandler(libvirtClient))
	mux.HandleFunc("DELETE /api/libvirt/vms/{name}", deleteVMHandler(libvirtClient))

	mux.HandleFunc("POST /api/libvirt/vms/{name}/nics", attachNICHandler(libvirtClient))
	mux.HandleFunc("GET /api/libvirt/vms/{name}/nics", listNICsHandler(libvirtClient))

	mux.HandleFunc("POST /api/ovn/switches", createSwitchHandler(ovnClient))
	mux.HandleFunc("GET /api/ovn/switches/{name}", getSwitchHandler(ovnClient))
	mux.HandleFunc("DELETE /api/ovn/switches/{name}", deleteSwitchHandler(ovnClient))

	mux.HandleFunc("POST /api/topology/nic-switch-connections", connectNICtoSwitchHandler(core))
	mux.HandleFunc("DELETE /api/topology/nic-switch-connections", disconnectNICfromSwitchHandler(core))

	mux.HandleFunc("POST /api/topology/switch-router-connections", connectSwitchToRouterHandler(core))
	mux.HandleFunc("DELETE /api/topology/switch-router-connections", disconnectSwitchFromRouterHandler(core))

	return mux
}
