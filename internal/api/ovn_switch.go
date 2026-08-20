package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
)

// CreateSwitchRequest is the request body for creating a irtual switch.
type CreateSwitchRequest struct {
	Name string `json:"name"`
}

// SwitchRespoinse is the ICP/IP-facing representation of a virtual switch.
// It intentionally hides OVN-specific vocabulary from the API surface.
type SwitchResponse struct {
	Name string `json:"name"`
}

// createSwitchHandler handles POST /api/ovn/switches
func createSwitchHandler(client *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSwitchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		swReq := &ovn.LogicalSwitch{
			Name: req.Name,
		}

		sw, err := client.CreateSwitch(r.Context(), swReq)
		if err != nil {
			if errors.Is(err, ovn.ErrSwitchAlreadyExists) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(toSwitchResponse(sw))
	}
}

// GetSwitch handles GET /api/ovn/switches/{name}.
func getSwitchHandler(client *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		sw, err := client.GetSwitch(r.Context(), name)
		if err != nil {
			if errors.Is(err, ovn.ErrSwitchNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toSwitchResponse(sw))
	}
}

// DeleteSwitch handles DELETE /api/ovn/switches/{name}.
func deleteSwitchHandler(client *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if err := client.DeleteSwitch(r.Context(), name); err != nil {
			if errors.Is(err, ovn.ErrSwitchNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func toSwitchResponse(sw *ovn.LogicalSwitch) SwitchResponse {
	return SwitchResponse{
		Name: sw.Name,
	}
}
