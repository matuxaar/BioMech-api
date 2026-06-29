package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/matuxaar/BioMech-api/internal/middleware"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/service"
)

var (
	activeConns   = make(map[*websocket.Conn]struct{})
	activeConnsMu sync.Mutex
)

func CloseAllWS() {
	activeConnsMu.Lock()
	defer activeConnsMu.Unlock()
	for conn := range activeConns {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutting down"))
		conn.Close()
	}
	activeConns = make(map[*websocket.Conn]struct{})
}

func trackConn(conn *websocket.Conn) {
	activeConnsMu.Lock()
	activeConns[conn] = struct{}{}
	activeConnsMu.Unlock()
}

func untrackConn(conn *websocket.Conn) {
	activeConnsMu.Lock()
	delete(activeConns, conn)
	activeConnsMu.Unlock()
}

type WSHandler struct {
	mlClient *service.MLClient
}

func NewWSHandler(mlClient *service.MLClient) *WSHandler {
	return &WSHandler{mlClient: mlClient}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	host := r.Header.Get("Host")
	return origin == "http://"+host || origin == "https://"+host
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func (h *WSHandler) PredictStream(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("ws upgrade error", "error", err)
		return
	}
	middleware.TrackWSOpen()
	trackConn(conn)
	defer untrackConn(conn)
	defer middleware.TrackWSClose()
	defer conn.Close()

	conn.SetReadLimit(65536)
	conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		return nil
	})

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	var mu sync.Mutex
	write := func(msg model.StreamPredictResponse) {
		mu.Lock()
		defer mu.Unlock()
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := conn.WriteJSON(msg); err != nil {
			slog.Warn("ws write error", "error", err)
		}
		middleware.TrackWSMessage("outgoing")
	}

	buffer := make([]model.StreamSample, 0, 256)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, raw, err := conn.ReadMessage()
		if err == nil {
			middleware.TrackWSMessage("incoming")
		}
		if err != nil {
			break
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

		var req model.StreamPredictRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			write(model.StreamPredictResponse{Type: "error", Error: "invalid request"})
			continue
		}

		if len(req.Samples) == 0 {
			continue
		}

		buffer = append(buffer, req.Samples...)

		if len(buffer) >= 32 {
			samples := make([]model.StreamSample, len(buffer))
			copy(samples, buffer)
			buffer = buffer[:0]

			go func(s []model.StreamSample) {
				select {
				case <-ctx.Done():
					return
				default:
				}
				predictions, err := h.mlClient.PredictStream(ctx, s)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					default:
					}
					write(model.StreamPredictResponse{Type: "error", Error: err.Error()})
					return
				}
				for _, p := range predictions {
					write(model.StreamPredictResponse{Type: "prediction", Prediction: p})
				}
			}(samples)
		}
	}
}
