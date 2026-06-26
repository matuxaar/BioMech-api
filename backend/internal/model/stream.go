package model

// StreamSample is a lightweight EMG sample for real-time WebSocket streaming.
type StreamSample struct {
	Channel1 float64 `json:"ch1"`
	Channel2 float64 `json:"ch2"`
	Channel3 float64 `json:"ch3"`
	Channel4 float64 `json:"ch4"`
	Channel5 float64 `json:"ch5"`
	Channel6 float64 `json:"ch6"`
	Channel7 float64 `json:"ch7"`
	Channel8 float64 `json:"ch8"`
}

// StreamPredictRequest is sent from the app to the backend via WebSocket.
type StreamPredictRequest struct {
	Samples  []StreamSample `json:"samples"`
	DeviceID string         `json:"device_id,omitempty"`
}

// StreamPredictResponse is sent from the backend to the app via WebSocket.
type StreamPredictResponse struct {
	Type       string `json:"type"`
	Prediction string `json:"prediction,omitempty"`
	Error      string `json:"error,omitempty"`
}
