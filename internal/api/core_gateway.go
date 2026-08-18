package api

import (
	"encoding/json"
	"net/http"

	"github.com/nakamuraitsuki/Girder/internal/core"
)

type CreateGatewayRequest struct {
	Name        string `json:"name"`
	PhysicalNICName string `json:"physicalNicName"`
}

func createGatewayHandler(core *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateGatewayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := core.CreateGateway(
			r.Context(),
			req.Name,
			req.PhysicalNICName,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
