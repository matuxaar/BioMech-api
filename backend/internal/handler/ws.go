package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/service"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	mlClient *service.MLClient
}

func NewWSHandler(mlClient *service.MLClient) *WSHandler {
	return &WSHandler{mlClient: mlClient}
}

func (h *WSHandler) PredictStream(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}
	defer conn.Close()

	var mu sync.Mutex
	write := func(msg model.StreamPredictResponse) {
		mu.Lock()
		defer mu.Unlock()
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("ws write error: %v", err)
		}
	}

	buffer := make([]model.StreamSample, 0, 256)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

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
				predictions, err := h.mlClient.PredictStream(s)
				if err != nil {
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
