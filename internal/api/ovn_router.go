package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
)

type createRouterRequest struct {
	Name string `json:"name"`
}

// createRouterHandler handles POST /api/ovn/routers
func createRouterHandler(ovnClient *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRouterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		router, err := ovnClient.CreateRouter(r.Context(), &ovn.LogicalRouter{
			Name: req.Name,
		})
		if err != nil {
			if errors.Is(err, ovn.ErrLogicalRouterAlreadyExists) {
				http.Error(w, "Router already exists", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to create router", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(router)
	}
}

// getRouterHandler handles GET /api/ovn/routers/{name}
func getRouterHandler(ovnClient *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			http.Error(w, "Missing router name", http.StatusBadRequest)
			return
		}

		router, err := ovnClient.GetRouter(r.Context(), name)
		if err != nil {
			if errors.Is(err, ovn.ErrLogicalRouterNotFound) {
				http.Error(w, "Router not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to get router", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(router)
	}
}

// deleteRouterHandler handles DELETE /api/ovn/routers/{name}
func deleteRouterHandler(ovnClient *ovn.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			http.Error(w, "Missing router name", http.StatusBadRequest)
			return
		}

		if err := ovnClient.DeleteRouter(r.Context(), name); err != nil {
			if errors.Is(err, ovn.ErrLogicalRouterNotFound) {
				http.Error(w, "Router not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to delete router", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
