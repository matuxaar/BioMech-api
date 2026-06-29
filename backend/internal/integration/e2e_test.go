//go:build e2e

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/handler"
	"github.com/matuxaar/BioMech-api/internal/middleware"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
	"github.com/matuxaar/BioMech-api/internal/service"
	"github.com/matuxaar/BioMech-api/internal/testhelper"
)

var testDB *testhelper.TestDB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	middleware.InitAuth(true)

	migrationsDir := findMigrationsDir()
	db, err := testhelper.StartDBForMain(migrationsDir)
	if err != nil {
		panic(err)
	}
	testDB = db
	code := m.Run()
	testDB.Pool.Close()
	if testDB.Cancel != nil {
		testDB.Cancel()
	}
	os.Exit(code)
}

func findMigrationsDir() string {
	candidates := []string{"../../migrations", "../../../migrations", "../migrations"}
	for _, c := range candidates {
		abs, _ := filepath.Abs(c)
		files, err := filepath.Glob(filepath.Join(abs, "*.sql"))
		if err == nil && len(files) > 0 {
			return abs
		}
	}
	dir, _ := os.Getwd()
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, "migrations")
		files, err := filepath.Glob(filepath.Join(candidate, "*.sql"))
		if err == nil && len(files) > 0 {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	panic("could not find migrations directory")
}

func setupE2EServer(t *testing.T) *gin.Engine {
	userRepo := repository.NewUserRepository(testDB.Pool)
	deviceRepo := repository.NewDeviceRepository(testDB.Pool)
	emgRepo := repository.NewEMGRepository(testDB.Pool)
	trainingRepo := repository.NewTrainingRepository(testDB.Pool)
	trainingFileRepo := repository.NewTrainingFileRepository(testDB.Pool)
	statsRepo := repository.NewStatsRepository(testDB.Pool)

	mlClient := service.NewMLClient("http://ml:8000", 0)
	statsService := service.NewStatsService(statsRepo)

	return handler.SetupRouter(
		&firebase.App{},
		handler.NewAuthHandler(service.NewAuthService(userRepo)),
		handler.NewUserHandler(service.NewAuthService(userRepo), "/tmp/avatars"),
		handler.NewDeviceHandler(service.NewDeviceService(deviceRepo)),
		handler.NewEMGHandler(service.NewEMGService(emgRepo, deviceRepo)),
		handler.NewTrainingHandler(service.NewTrainingService(trainingRepo, emgRepo, deviceRepo, mlClient)),
		handler.NewStatsHandler(statsService),
		handler.NewWSHandler(mlClient),
		handler.NewTrainingFileHandler(service.NewTrainingFileService(trainingFileRepo, "/tmp/training")),
		50, "/tmp/uploads",
	)
}

func TestE2EHealthEndpoint(t *testing.T) {
	router := setupE2EServer(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %s", body["status"])
	}
}

func TestE2EMetricsEndpoint(t *testing.T) {
	router := setupE2EServer(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestE2EUserLifecycle(t *testing.T) {
	router := setupE2EServer(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/firebase", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestE2EDeviceCRUD(t *testing.T) {
	router := setupE2EServer(t)

	body := map[string]any{
		"name":       "E2E Test Band",
		"type":       "prosthetic",
		"hw_version": "2.0",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices",
		bytes.NewReader(mustJSON(t, body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var device model.Device
	json.Unmarshal(w.Body.Bytes(), &device)
	if device.Name != "E2E Test Band" {
		t.Errorf("expected E2E Test Band, got %s", device.Name)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/devices", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/devices/"+device.ID, nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for GetByID, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/devices/"+device.ID, nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for Delete, got %d", w.Code)
	}
}

func TestE2EEMGSession(t *testing.T) {
	router := setupE2EServer(t)

	devBody := map[string]any{
		"name": "EMG Band", "type": "sensor", "hw_version": "1.0",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/devices",
		bytes.NewReader(mustJSON(t, devBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var device model.Device
	json.Unmarshal(w.Body.Bytes(), &device)

	sessBody := map[string]any{
		"device_id": device.ID,
		"label":     "e2e-test",
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/emg/sessions",
		bytes.NewReader(mustJSON(t, sessBody)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var session model.EMGSession
	json.Unmarshal(w.Body.Bytes(), &session)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/emg/sessions", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func mustJSON(t *testing.T, v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
		return nil
	}
	return b
}
