package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/matuxaar/BioMech-api/internal/model"
)

func TestMLClient_Train_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/train" {
			t.Errorf("expected /train, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":     "completed",
			"accuracy":   0.95,
			"model_path": "/app/models/test-job.h5",
		})
	}))
	defer srv.Close()

	client := NewMLClient(srv.URL, 5*time.Second)
	job := &model.TrainingJob{
		ID:         "test-job",
		SessionIDs: []string{"sess-1", "sess-2"},
	}

	result, err := client.Train(context.Background(), job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "completed" {
		t.Errorf("expected status completed, got %s", result.Status)
	}
	if result.Accuracy != 0.95 {
		t.Errorf("expected accuracy 0.95, got %f", result.Accuracy)
	}
	if result.ModelPath != "/app/models/test-job.h5" {
		t.Errorf("expected model path /app/models/test-job.h5, got %s", result.ModelPath)
	}
}

func TestMLClient_Train_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewMLClient(srv.URL, 5*time.Second)
	job := &model.TrainingJob{
		ID:         "test-job",
		SessionIDs: []string{"sess-1"},
	}

	_, err := client.Train(context.Background(), job)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMLClient_Train_UnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json`))
	}))
	defer srv.Close()

	client := NewMLClient(srv.URL, 5*time.Second)
	job := &model.TrainingJob{
		ID:         "test-job",
		SessionIDs: []string{"sess-1"},
	}

	_, err := client.Train(context.Background(), job)
	if err == nil {
		t.Fatal("expected unmarshal error, got nil")
	}
}

func TestMLClient_Predict_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/predict/stream" {
			t.Errorf("expected /predict/stream, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"predictions": []string{"fist", "open", "rest"},
		})
	}))
	defer srv.Close()

	client := NewMLClient(srv.URL, 5*time.Second)
	samples := []model.EMGSample{
		{Channel1: 0.1, Channel2: 0.2, Channel3: 0.3, Channel4: 0.4, Channel5: 0.5, Channel6: 0.6, Channel7: 0.7, Channel8: 0.8},
		{Channel1: 0.2, Channel2: 0.3, Channel3: 0.4, Channel4: 0.5, Channel5: 0.6, Channel6: 0.7, Channel7: 0.8, Channel8: 0.9},
	}

	predictions, err := client.Predict(context.Background(), samples)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(predictions) != 3 {
		t.Fatalf("expected 3 predictions, got %d", len(predictions))
	}
	if predictions[0] != "fist" {
		t.Errorf("expected fist, got %s", predictions[0])
	}
}

func TestMLClient_PredictStream_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/predict/stream" {
			t.Errorf("expected /predict/stream, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"predictions": []string{"point"},
		})
	}))
	defer srv.Close()

	client := NewMLClient(srv.URL, 5*time.Second)
	samples := []model.StreamSample{
		{Channel1: 0.1, Channel2: 0.2, Channel3: 0.3, Channel4: 0.4, Channel5: 0.5, Channel6: 0.6, Channel7: 0.7, Channel8: 0.8},
	}

	predictions, err := client.PredictStream(context.Background(), samples)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(predictions) != 1 {
		t.Fatalf("expected 1 prediction, got %d", len(predictions))
	}
	if predictions[0] != "point" {
		t.Errorf("expected point, got %s", predictions[0])
	}
}

func TestMLClient_Predict_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client := NewMLClient(srv.URL, 5*time.Second)
	_, err := client.Predict(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
