package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
)

type setRouterInterfaceAddressRequest struct {
	Router    string `json:"router"`
	Interface string `json:"interface"`
	Address   string `json:"address"`
}

func setRouterInterfaceAddressHandler(ovnClient *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req setRouterInterfaceAddressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Router == "" || req.Interface == "" || req.Address == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		router, err := ovnClient.GetRouter(r.Context(), req.Router)
		if err != nil {
			if errors.Is(err, ovn.ErrLogicalRouterNotFound) {
				http.Error(w, "Router not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to get router", http.StatusInternalServerError)
			return
		}

		if err := ovnClient.SetLogicalRouterPortNetworks(
			r.Context(),
			router,
			req.Interface,
			[]string{req.Address},
		); err != nil {
			if errors.Is(err, ovn.ErrLogicalRouterPortNotFound) {
				http.Error(w, "Router interface not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to set router interface address", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type SetDefaultRouteRequest struct {
	Nexthop string `json:"nexthop"`
}

// POST /api/ovn/routers/{name}/default-route
func setDefaultRouteHandler(ovnClient *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routerName := r.PathValue("name")

		var req SetDefaultRouteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Nexthop == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		router, err := ovnClient.GetRouter(r.Context(), routerName)
		if err != nil {
			http.Error(w, "Failed to get router", http.StatusNotFound)
			return
		}

		route, err := ovnClient.ReplaceStaticRoute(
			r.Context(),
			router,
			&ovn.LogicalRouterStaticRoute{
				IPPrefix: "0.0.0.0/0",
				Nexthop:  req.Nexthop,
			},
		)
		if err != nil {
			http.Error(w, "Failed to set default route", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(route)
	}
}