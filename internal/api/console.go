package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nakamuraitsuki/Girder/infrastructure/libvirt"
)

var consoleUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func consoleHandler(client *libvirt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		conn, err := consoleUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		stream, err := client.OpenConsole(name)
		if err != nil {
			_ = conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseInternalServerErr,
					err.Error(),
				),
			)
			return
		}
		defer stream.Free()

		errCh := make(chan error, 2)

		// WebSocket -> libvirt console
		go func() {
			for {
				messageType, data, err := conn.ReadMessage()
				if err != nil {
					errCh <- err
					return
				}

				if messageType != websocket.TextMessage &&
					messageType != websocket.BinaryMessage {
					continue
				}

				if _, err := stream.Send(data); err != nil {
					errCh <- fmt.Errorf("send to console: %w", err)
					return
				}
			}
		}()

		// libvirt console -> WebSocket
		go func() {
			buf := make([]byte, 4096)

			for {
				n, err := stream.Recv(buf)
				if err != nil {
					errCh <- fmt.Errorf("receive from console: %w", err)
					return
				}

				if n == 0 {
					continue
				}

				if err := conn.WriteMessage(
					websocket.BinaryMessage,
					buf[:n],
				); err != nil {
					errCh <- fmt.Errorf("write websocket: %w", err)
					return
				}
			}
		}()

		<-errCh

		stream.Abort()
	}
}