package api

import (
	"encoding/json"
	"net/http"

	"github.com/nakamuraitsuki/Girder/internal/core"
)

type switchRouterRequest struct {
	SwitchName string `json:"switchName"`
	RouterName string `json:"routerName"`
}

func connectSwitchToRouterHandler(c *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req switchRouterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.SwitchName == "" || req.RouterName == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		err := c.ConnectSwitchToRouter(r.Context(), req.SwitchName, req.RouterName)
		if err != nil {
			http.Error(w, "Failed to connect switch to router", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func disconnectSwitchFromRouterHandler(c *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req switchRouterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.SwitchName == "" || req.RouterName == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		err := c.DisconnectSwitchFromRouter(r.Context(), req.SwitchName, req.RouterName)
		if err != nil {
			http.Error(w, "Failed to disconnect switch from router", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}