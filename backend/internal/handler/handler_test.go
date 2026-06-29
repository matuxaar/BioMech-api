package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/service"
)

// --- Mock implementations ---

type mockAuthService struct {
	syncUserFn      func(ctx context.Context, firebaseUID, email string) (*model.User, error)
	getProfileFn    func(ctx context.Context, userID string) (*model.ProfileResponse, error)
	updateProfileFn func(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error)
}

func (m *mockAuthService) SyncUser(ctx context.Context, firebaseUID, email string) (*model.User, error) {
	return m.syncUserFn(ctx, firebaseUID, email)
}
func (m *mockAuthService) GetProfile(ctx context.Context, userID string) (*model.ProfileResponse, error) {
	return m.getProfileFn(ctx, userID)
}
func (m *mockAuthService) UpdateProfile(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
	return m.updateProfileFn(ctx, userID, req)
}

type mockDeviceService struct {
	createFn     func(ctx context.Context, userID string, req *model.CreateDeviceRequest) (*model.Device, error)
	listByUserFn func(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.Device], error)
	getByIDFn    func(ctx context.Context, userID, id string) (*model.Device, error)
	updateFn     func(ctx context.Context, userID, id string, req *model.UpdateDeviceRequest) (*model.Device, error)
	deleteFn     func(ctx context.Context, userID, id string) error
	getActionsFn func(ctx context.Context, userID, id string) (*model.DeviceActionsResponse, error)
}

func (m *mockDeviceService) Create(ctx context.Context, userID string, req *model.CreateDeviceRequest) (*model.Device, error) {
	return m.createFn(ctx, userID, req)
}
func (m *mockDeviceService) ListByUser(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.Device], error) {
	return m.listByUserFn(ctx, userID, page, limit)
}
func (m *mockDeviceService) GetByID(ctx context.Context, userID, id string) (*model.Device, error) {
	return m.getByIDFn(ctx, userID, id)
}
func (m *mockDeviceService) Update(ctx context.Context, userID, id string, req *model.UpdateDeviceRequest) (*model.Device, error) {
	return m.updateFn(ctx, userID, id, req)
}
func (m *mockDeviceService) Delete(ctx context.Context, userID, id string) error {
	return m.deleteFn(ctx, userID, id)
}
func (m *mockDeviceService) GetActions(ctx context.Context, userID, id string) (*model.DeviceActionsResponse, error) {
	return m.getActionsFn(ctx, userID, id)
}

type mockEMGService struct {
	startSessionFn    func(ctx context.Context, userID string, req *model.CreateEMGSessionRequest) (*model.EMGSession, error)
	endSessionFn      func(ctx context.Context, userID, id string) error
	listSessionsFn    func(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.EMGSession], error)
	getSessionFn      func(ctx context.Context, userID, id string) (*model.EMGSession, error)
	addSampleFn       func(ctx context.Context, userID, sessionID string, req *model.AddSampleRequest) (*model.EMGSample, error)
	addSamplesBatchFn func(ctx context.Context, userID, sessionID string, samples []model.AddSampleRequest) error
	getSamplesFn      func(ctx context.Context, userID, sessionID string, page, limit int) (*model.PaginatedResponse[model.EMGSample], error)
}

func (m *mockEMGService) StartSession(ctx context.Context, userID string, req *model.CreateEMGSessionRequest) (*model.EMGSession, error) {
	return m.startSessionFn(ctx, userID, req)
}
func (m *mockEMGService) EndSession(ctx context.Context, userID, id string) error {
	return m.endSessionFn(ctx, userID, id)
}
func (m *mockEMGService) ListSessions(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.EMGSession], error) {
	return m.listSessionsFn(ctx, userID, page, limit)
}
func (m *mockEMGService) GetSession(ctx context.Context, userID, id string) (*model.EMGSession, error) {
	return m.getSessionFn(ctx, userID, id)
}
func (m *mockEMGService) AddSample(ctx context.Context, userID, sessionID string, req *model.AddSampleRequest) (*model.EMGSample, error) {
	return m.addSampleFn(ctx, userID, sessionID, req)
}
func (m *mockEMGService) AddSamplesBatch(ctx context.Context, userID, sessionID string, samples []model.AddSampleRequest) error {
	return m.addSamplesBatchFn(ctx, userID, sessionID, samples)
}
func (m *mockEMGService) GetSamples(ctx context.Context, userID, sessionID string, page, limit int) (*model.PaginatedResponse[model.EMGSample], error) {
	return m.getSamplesFn(ctx, userID, sessionID, page, limit)
}

type mockTrainingService struct {
	createJobFn       func(ctx context.Context, userID string, req *model.CreateTrainingJobRequest) (*model.TrainingJob, error)
	startTrainingFn   func(ctx context.Context, jobID string) error
	listJobsFn        func(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingJob], error)
	getJobFn          func(ctx context.Context, userID, id string) (*model.TrainingJob, error)
	predictFn         func(ctx context.Context, samples []model.EMGSample) (*model.PredictResponse, error)
	processUploadFn   func(ctx context.Context, userID, deviceID, label string, file io.Reader) (*model.EMGSession, error)
	updateJobStatusFn func(ctx context.Context, id, status, modelPath string, accuracy float64, errMsg string) error
}

func (m *mockTrainingService) CreateJob(ctx context.Context, userID string, req *model.CreateTrainingJobRequest) (*model.TrainingJob, error) {
	return m.createJobFn(ctx, userID, req)
}
func (m *mockTrainingService) StartTraining(ctx context.Context, jobID string) error {
	return m.startTrainingFn(ctx, jobID)
}
func (m *mockTrainingService) ListJobs(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingJob], error) {
	return m.listJobsFn(ctx, userID, page, limit)
}
func (m *mockTrainingService) GetJob(ctx context.Context, userID, id string) (*model.TrainingJob, error) {
	return m.getJobFn(ctx, userID, id)
}
func (m *mockTrainingService) Predict(ctx context.Context, samples []model.EMGSample) (*model.PredictResponse, error) {
	return m.predictFn(ctx, samples)
}
func (m *mockTrainingService) ProcessUpload(ctx context.Context, userID, deviceID, label string, file io.Reader) (*model.EMGSession, error) {
	return m.processUploadFn(ctx, userID, deviceID, label, file)
}
func (m *mockTrainingService) UpdateJobStatus(ctx context.Context, id, status, modelPath string, accuracy float64, errMsg string) error {
	return m.updateJobStatusFn(ctx, id, status, modelPath, accuracy, errMsg)
}

type mockStatsService struct {
	getDashboardStatsFn func(ctx context.Context, userID string) (*model.DashboardStats, error)
}

func (m *mockStatsService) GetDashboardStats(ctx context.Context, userID string) (*model.DashboardStats, error) {
	return m.getDashboardStatsFn(ctx, userID)
}

type mockTrainingFileService struct {
	uploadFn func(ctx context.Context, userID, deviceID, label, filename string, file io.Reader, size int64) (*model.TrainingFile, error)
	listFn   func(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingFile], error)
	getFn    func(ctx context.Context, id, userID string) (*model.TrainingFile, error)
	deleteFn func(ctx context.Context, id, userID string) error
}

func (m *mockTrainingFileService) Upload(ctx context.Context, userID, deviceID, label, filename string, file io.Reader, size int64) (*model.TrainingFile, error) {
	return m.uploadFn(ctx, userID, deviceID, label, filename, file, size)
}
func (m *mockTrainingFileService) List(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingFile], error) {
	return m.listFn(ctx, userID, page, limit)
}
func (m *mockTrainingFileService) Get(ctx context.Context, id, userID string) (*model.TrainingFile, error) {
	return m.getFn(ctx, id, userID)
}
func (m *mockTrainingFileService) Delete(ctx context.Context, id, userID string) error {
	return m.deleteFn(ctx, id, userID)
}

// --- Test helpers ---

func setupTestRouter(h ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	for _, m := range h {
		r.Use(m)
	}
	return r
}

func authMiddleware(userID, email string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("email", email)
		c.Next()
	}
}

func jsonBody(v any) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

// --- Auth handler tests ---

func TestAuthHandler_SyncUser_Success(t *testing.T) {
	mock := &mockAuthService{
		syncUserFn: func(ctx context.Context, uid, email string) (*model.User, error) {
			return &model.User{ID: uid, Email: email}, nil
		},
	}
	h := NewAuthHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.SyncUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthHandler_SyncUser_ServiceError(t *testing.T) {
	mock := &mockAuthService{
		syncUserFn: func(ctx context.Context, uid, email string) (*model.User, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewAuthHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.SyncUser)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- User handler tests ---

func TestUserHandler_Me_Success(t *testing.T) {
	mock := &mockAuthService{
		getProfileFn: func(ctx context.Context, userID string) (*model.ProfileResponse, error) {
			return &model.ProfileResponse{ID: userID, Email: "test@test.com"}, nil
		},
	}
	h := NewUserHandler(mock, "/avatars")

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Me)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp model.ProfileResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != "user-1" {
		t.Errorf("expected user-1, got %s", resp.ID)
	}
}

func TestUserHandler_Me_NotFound(t *testing.T) {
	mock := &mockAuthService{
		getProfileFn: func(ctx context.Context, userID string) (*model.ProfileResponse, error) {
			return nil, errors.New("user not found")
		},
	}
	h := NewUserHandler(mock, "/avatars")

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Me)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUserHandler_Update_ValidationError(t *testing.T) {
	mock := &mockAuthService{
		updateProfileFn: func(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
			return &model.User{ID: userID}, nil
		},
		getProfileFn: func(ctx context.Context, userID string) (*model.ProfileResponse, error) {
			return &model.ProfileResponse{ID: userID}, nil
		},
	}
	h := NewUserHandler(mock, "/avatars")

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Update)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Update_NicknameTaken(t *testing.T) {
	mock := &mockAuthService{
		updateProfileFn: func(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
			return nil, service.ErrNicknameTaken
		},
	}
	h := NewUserHandler(mock, "/avatars")

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Update)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", jsonBody(model.UpdateUserRequest{}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

// --- Device handler tests ---

func TestDeviceHandler_Create_Success(t *testing.T) {
	mock := &mockDeviceService{
		createFn: func(ctx context.Context, userID string, req *model.CreateDeviceRequest) (*model.Device, error) {
			return &model.Device{ID: "dev-1", Name: req.Name, UserID: userID}, nil
		},
	}
	h := NewDeviceHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", jsonBody(model.CreateDeviceRequest{
		Name:     "My Band",
		Type:     "prosthetic",
		HWVersion: "1.0",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestDeviceHandler_Create_BadRequest(t *testing.T) {
	h := NewDeviceHandler(&mockDeviceService{})

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeviceHandler_GetByID_NotFound(t *testing.T) {
	mock := &mockDeviceService{
		getByIDFn: func(ctx context.Context, userID, id string) (*model.Device, error) {
			return nil, errors.New("device not found")
		},
	}
	h := NewDeviceHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"))
	r.GET("/:id", h.GetByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeviceHandler_Delete_AccessDenied(t *testing.T) {
	mock := &mockDeviceService{
		deleteFn: func(ctx context.Context, userID, id string) error {
			return service.ErrAccessDenied
		},
	}
	h := NewDeviceHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"))
	r.DELETE("/:id", h.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/dev-1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

// --- EMG handler tests ---

func TestEMGHandler_StartSession_Success(t *testing.T) {
	mock := &mockEMGService{
		startSessionFn: func(ctx context.Context, userID string, req *model.CreateEMGSessionRequest) (*model.EMGSession, error) {
			return &model.EMGSession{ID: "sess-1", UserID: userID, DeviceID: req.DeviceID}, nil
		},
	}
	h := NewEMGHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.StartSession)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", jsonBody(model.CreateEMGSessionRequest{DeviceID: "dev-1"}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestEMGHandler_AddSamplesBatch_BadRequest(t *testing.T) {
	h := NewEMGHandler(&mockEMGService{})

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"))
	r.POST("/:id", h.AddSamplesBatch)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/sess-1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Training handler tests ---

func TestTrainingHandler_CreateJob_Success(t *testing.T) {
	mock := &mockTrainingService{
		createJobFn: func(ctx context.Context, userID string, req *model.CreateTrainingJobRequest) (*model.TrainingJob, error) {
			return &model.TrainingJob{ID: "job-1", UserID: userID}, nil
		},
		startTrainingFn: func(ctx context.Context, jobID string) error {
			return nil
		},
	}
	h := NewTrainingHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.CreateJob)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", jsonBody(model.CreateTrainingJobRequest{SessionIDs: []string{"sess-1"}}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestTrainingHandler_GetJob_NotFound(t *testing.T) {
	mock := &mockTrainingService{
		getJobFn: func(ctx context.Context, userID, id string) (*model.TrainingJob, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewTrainingHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"))
	r.GET("/:id", h.GetJob)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestTrainingHandler_Predict_Success(t *testing.T) {
	mock := &mockTrainingService{
		predictFn: func(ctx context.Context, samples []model.EMGSample) (*model.PredictResponse, error) {
			return &model.PredictResponse{Predictions: []string{"fist"}}, nil
		},
	}
	h := NewTrainingHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Predict)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", jsonBody(predictRequest{
		Samples: []model.AddSampleRequest{{Channel1: 0.5}},
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTrainingHandler_Upload_BadRequest(t *testing.T) {
	mock := &mockTrainingService{
		processUploadFn: func(ctx context.Context, userID, deviceID, label string, file io.Reader) (*model.EMGSession, error) {
			return nil, nil
		},
	}
	h := NewTrainingHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Upload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTrainingHandler_UpdateJobStatus_Success(t *testing.T) {
	mock := &mockTrainingService{
		updateJobStatusFn: func(ctx context.Context, id, status, modelPath string, accuracy float64, errMsg string) error {
			if id != "job-1" {
				t.Errorf("expected job-1, got %s", id)
			}
			return nil
		},
	}
	h := NewTrainingHandler(mock)

	r := setupTestRouter()
	r.POST("/:id", h.UpdateJobStatus)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/job-1", jsonBody(map[string]any{
		"status": "completed",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// --- Stats handler tests ---

func TestStatsHandler_Dashboard_Success(t *testing.T) {
	mock := &mockStatsService{
		getDashboardStatsFn: func(ctx context.Context, userID string) (*model.DashboardStats, error) {
			return &model.DashboardStats{TotalTrainings: 10}, nil
		},
	}
	h := NewStatsHandler(mock)

	r := setupTestRouter(authMiddleware("user-1", "test@test.com"), h.Dashboard)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}


