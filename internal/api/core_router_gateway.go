package api

import (
	"encoding/json"
	"net/http"

	"github.com/nakamuraitsuki/Girder/internal/core"
)

type routerGatewayRequest struct {
	RouterName  string `json:"routerName"`
	GatewayName string `json:"gatewayName"`
}

func connectRouterToGatewayHandler(c *core.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req routerGatewayRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if req.RouterName == "" || req.GatewayName == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if err := c.ConnectSwitchToRouter(
			r.Context(),
			req.GatewayName,
			req.RouterName,
		); err != nil {
			http.Error(w, "Failed to connect router to gateway", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
