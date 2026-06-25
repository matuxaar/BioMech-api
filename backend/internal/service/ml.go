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

type predictRequest struct {
	Samples []model.EMGSample `json:"samples"`
}

type predictResponse struct {
	Predictions []string `json:"predictions"`
}

func (c *MLClient) Predict(samples []model.EMGSample) ([]string, error) {
	body := predictRequest{Samples: samples}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(c.baseURL+"/predict", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result predictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Predictions, nil
}
