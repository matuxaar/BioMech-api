package service

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/matuxaar/BioMech-api/internal/model"
)

type MLClient struct {
	baseURL string
	client  *http.Client
}

func NewMLClient(baseURL string) *MLClient {
	return &MLClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type trainRequest struct {
	JobID      string   `json:"job_id"`
	SessionIDs []string `json:"session_ids"`
}

func (c *MLClient) Train(job *model.TrainingJob) error {
	body := trainRequest{
		JobID:      job.ID,
		SessionIDs: job.SessionIDs,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.baseURL+"/train", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type mlPredictRequest struct {
	Samples []mlStreamSample `json:"samples"`
}

type mlStreamSample struct {
	Channel1 float64 `json:"channel_1"`
	Channel2 float64 `json:"channel_2"`
	Channel3 float64 `json:"channel_3"`
	Channel4 float64 `json:"channel_4"`
	Channel5 float64 `json:"channel_5"`
	Channel6 float64 `json:"channel_6"`
	Channel7 float64 `json:"channel_7"`
	Channel8 float64 `json:"channel_8"`
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
	resp, err := c.client.Post(c.baseURL+"/predict", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
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

	var result mlPredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Predictions, nil
}
