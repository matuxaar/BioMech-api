package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/matuxaar/BioMech-api/internal/model"
)

type MLClient struct {
	baseURL string
	client  *http.Client
}

func NewMLClient(baseURL string) *MLClient {
	return &MLClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 300 {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("ML service returned %d: %s", resp.StatusCode, string(body[:n]))
	}
	return nil
}

type trainRequest struct {
	JobID      string   `json:"job_id"`
	SessionIDs []string `json:"session_ids"`
}

type trainResponse struct {
	Status    string  `json:"status"`
	Accuracy  float64 `json:"accuracy"`
	ModelPath string  `json:"model_path"`
}

func (c *MLClient) Train(job *model.TrainingJob) (*trainResponse, error) {
	body := trainRequest{
		JobID:      job.ID,
		SessionIDs: job.SessionIDs,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/train", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey := os.Getenv("INTERNAL_API_KEY"); apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var result trainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type mlPredictRequest struct {
	Samples []mlStreamSample `json:"samples"`
}

type mlStreamSample struct {
	Channel1 float64 `json:"ch1"`
	Channel2 float64 `json:"ch2"`
	Channel3 float64 `json:"ch3"`
	Channel4 float64 `json:"ch4"`
	Channel5 float64 `json:"ch5"`
	Channel6 float64 `json:"ch6"`
	Channel7 float64 `json:"ch7"`
	Channel8 float64 `json:"ch8"`
}

type mlPredictResponse struct {
	Predictions []string `json:"predictions"`
}

func (c *MLClient) Predict(samples []model.EMGSample) ([]string, error) {
	body := mlPredictRequest{Samples: make([]mlStreamSample, len(samples))}
	for i, s := range samples {
		body.Samples[i] = mlStreamSample{
			Channel1: s.Channel1,
			Channel2: s.Channel2,
			Channel3: s.Channel3,
			Channel4: s.Channel4,
			Channel5: s.Channel5,
			Channel6: s.Channel6,
			Channel7: s.Channel7,
			Channel8: s.Channel8,
		}
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Post(c.baseURL+"/predict/stream", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var result mlPredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Predictions, nil
}

func (c *MLClient) PredictStream(samples []model.StreamSample) ([]string, error) {
	converted := make([]mlStreamSample, len(samples))
	for i, s := range samples {
		converted[i] = mlStreamSample{
			Channel1: s.Channel1,
			Channel2: s.Channel2,
			Channel3: s.Channel3,
			Channel4: s.Channel4,
			Channel5: s.Channel5,
			Channel6: s.Channel6,
			Channel7: s.Channel7,
			Channel8: s.Channel8,
		}
	}

	body := mlPredictRequest{Samples: converted}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(c.baseURL+"/predict/stream", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var result mlPredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Predictions, nil
}
