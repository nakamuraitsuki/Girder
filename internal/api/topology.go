package api

import (
	"encoding/json"
	"net/http"

	"github.com/nakamuraitsuki/Girder/internal/core"
)

type nicSwitchRequest struct {
	VMName         string `json:"vmName"`
	NICLogicalName string `json:"nicLogicalName"`
	SwitchName     string `json:"switchName"`
}

func connectNICtoSwitchHandler(c *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req nicSwitchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.VMName == "" || req.NICLogicalName == "" || req.SwitchName == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		err := c.ConnectNICtoSwitch(r.Context(), req.VMName, req.NICLogicalName, req.SwitchName)
		if err != nil {
			http.Error(w, "Failed to connect NIC to switch", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func disconnectNICfromSwitchHandler(c *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req nicSwitchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.VMName == "" || req.NICLogicalName == "" || req.SwitchName == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		err := c.DisconnectNICfromSwitch(r.Context(), req.VMName, req.NICLogicalName, req.SwitchName)
		if err != nil {
			http.Error(w, "Failed to disconnect NIC from switch", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
