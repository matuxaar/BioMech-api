package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/testhelper"
)

var testDB *testhelper.TestDB

func TestMain(m *testing.M) {
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

func TestUserRepository_CreateAndFind(t *testing.T) {
	repo := NewUserRepository(testDB.Pool)
	ctx := context.Background()

	user, err := repo.Create(ctx, "test-user-1", "test@test.com")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if user.ID != "test-user-1" {
		t.Errorf("expected test-user-1, got %s", user.ID)
	}

	found, err := repo.FindByID(ctx, "test-user-1")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Email != "test@test.com" {
		t.Errorf("expected test@test.com, got %s", found.Email)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	repo := NewUserRepository(testDB.Pool)
	_, err := repo.FindByID(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestUserRepository_NicknameUniqueness(t *testing.T) {
	repo := NewUserRepository(testDB.Pool)
	ctx := context.Background()

	repo.Create(ctx, "user-a", "a@test.com")
	repo.Create(ctx, "user-b", "b@test.com")

	nickname := "shared-nick"
	_, err := repo.UpdateProfile(ctx, "user-a", &model.UpdateUserRequest{Nickname: &nickname})
	if err != nil {
		t.Fatalf("UpdateProfile failed: %v", err)
	}

	found, err := repo.FindByNickname(ctx, "shared-nick")
	if err != nil {
		t.Fatalf("FindByNickname failed: %v", err)
	}
	if found.ID != "user-a" {
		t.Errorf("expected user-a, got %s", found.ID)
	}
}

func TestUserRepository_CountDevices(t *testing.T) {
	userRepo := NewUserRepository(testDB.Pool)
	deviceRepo := NewDeviceRepository(testDB.Pool)
	ctx := context.Background()

	userRepo.Create(ctx, "user-stats", "stats@test.com")
	deviceRepo.Create(ctx, "user-stats", &model.CreateDeviceRequest{
		Name: "Band", Type: "prosthetic", HWVersion: "1.0",
	})
	deviceRepo.Create(ctx, "user-stats", &model.CreateDeviceRequest{
		Name: "Band2", Type: "sensor", HWVersion: "2.0",
	})

	count, err := userRepo.CountDevices(ctx, "user-stats")
	if err != nil {
		t.Fatalf("CountDevices failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 devices, got %d", count)
	}
}

func TestDeviceRepository_CRUD(t *testing.T) {
	userRepo := NewUserRepository(testDB.Pool)
	deviceRepo := NewDeviceRepository(testDB.Pool)
	ctx := context.Background()

	userRepo.Create(ctx, "dev-user", "dev@test.com")

	created, err := deviceRepo.Create(ctx, "dev-user", &model.CreateDeviceRequest{
		Name: "Test Band", Type: "prosthetic", HWVersion: "1.0",
		BLEServiceUUID: "abc-123",
	})
	if err != nil {
		t.Fatalf("Create device failed: %v", err)
	}
	if created.Name != "Test Band" {
		t.Errorf("expected Test Band, got %s", created.Name)
	}

	found, err := deviceRepo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.UserID != "dev-user" {
		t.Errorf("expected dev-user, got %s", found.UserID)
	}

	count, err := deviceRepo.CountByUserID(ctx, "dev-user")
	if err != nil {
		t.Fatalf("CountByUserID failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 device, got %d", count)
	}

	newName := "Updated Band"
	updated, err := deviceRepo.Update(ctx, created.ID, &model.UpdateDeviceRequest{Name: &newName})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != "Updated Band" {
		t.Errorf("expected Updated Band, got %s", updated.Name)
	}

	if err := deviceRepo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = deviceRepo.FindByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestTrainingRepository_CRUD(t *testing.T) {
	userRepo := NewUserRepository(testDB.Pool)
	deviceRepo := NewDeviceRepository(testDB.Pool)
	emgRepo := NewEMGRepository(testDB.Pool)
	trainingRepo := NewTrainingRepository(testDB.Pool)
	ctx := context.Background()

	userRepo.Create(ctx, "training-user", "train@test.com")
	device, _ := deviceRepo.Create(ctx, "training-user", &model.CreateDeviceRequest{
		Name: "Band", Type: "sensor", HWVersion: "1.0",
	})
	sess, _ := emgRepo.CreateSession(ctx, "training-user", &model.CreateEMGSessionRequest{
		DeviceID: device.ID, Label: "train-sess",
	})

	job, err := trainingRepo.Create(ctx, "training-user", &model.CreateTrainingJobRequest{
		SessionIDs: []string{sess.ID},
	})
	if err != nil {
		t.Fatalf("Create training job failed: %v", err)
	}
	if job.UserID != "training-user" {
		t.Errorf("expected training-user, got %s", job.UserID)
	}
	if job.Status != model.TrainingStatusPending {
		t.Errorf("expected pending, got %s", job.Status)
	}

	found, err := trainingRepo.FindByID(ctx, job.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.ID != job.ID {
		t.Errorf("expected %s, got %s", job.ID, found.ID)
	}

	err = trainingRepo.UpdateStatus(ctx, job.ID, model.TrainingStatusCompleted, "/models/test.h5", 0.95, "")
	if err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	updated, err := trainingRepo.FindByID(ctx, job.ID)
	if err != nil {
		t.Fatalf("FindByID after update failed: %v", err)
	}
	if updated.Status != model.TrainingStatusCompleted {
		t.Errorf("expected completed, got %s", updated.Status)
	}
	if updated.Accuracy != 0.95 {
		t.Errorf("expected 0.95, got %f", updated.Accuracy)
	}
}

func TestEMGRepository_SessionLifecycle(t *testing.T) {
	userRepo := NewUserRepository(testDB.Pool)
	deviceRepo := NewDeviceRepository(testDB.Pool)
	emgRepo := NewEMGRepository(testDB.Pool)
	ctx := context.Background()

	userRepo.Create(ctx, "emg-user", "emg@test.com")
	device, _ := deviceRepo.Create(ctx, "emg-user", &model.CreateDeviceRequest{
		Name: "EMG Band", Type: "sensor", HWVersion: "1.0",
	})

	session, err := emgRepo.CreateSession(ctx, "emg-user", &model.CreateEMGSessionRequest{
		DeviceID: device.ID, Label: "test-session",
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if session.Label != "test-session" {
		t.Errorf("expected test-session, got %s", session.Label)
	}

	err = emgRepo.EndSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("EndSession failed: %v", err)
	}

	ended, err := emgRepo.FindSessionByID(ctx, session.ID)
	if err != nil {
		t.Fatalf("FindSessionByID failed: %v", err)
	}
	if ended.EndedAt == nil {
		t.Error("expected ended_at to be set")
	}
}

func TestEMGRepository_Samples(t *testing.T) {
	userRepo := NewUserRepository(testDB.Pool)
	deviceRepo := NewDeviceRepository(testDB.Pool)
	emgRepo := NewEMGRepository(testDB.Pool)
	ctx := context.Background()

	userRepo.Create(ctx, "sample-user", "s@test.com")
	device, _ := deviceRepo.Create(ctx, "sample-user", &model.CreateDeviceRequest{
		Name: "Band", Type: "prosthetic", HWVersion: "1.0",
	})
	session, _ := emgRepo.CreateSession(ctx, "sample-user", &model.CreateEMGSessionRequest{DeviceID: device.ID})

	sample, err := emgRepo.AddSample(ctx, session.ID, &model.AddSampleRequest{
		Channel1: 0.5, Channel2: 0.1, Channel3: 0.2, Channel4: 0.3,
		Channel5: 0.4, Channel6: 0.5, Channel7: 0.6, Channel8: 0.7,
	})
	if err != nil {
		t.Fatalf("AddSample failed: %v", err)
	}
	if sample.Channel1 != 0.5 {
		t.Errorf("expected channel_1=0.5, got %f", sample.Channel1)
	}

	batchSamples := []model.AddSampleRequest{
		{Channel1: 1.0}, {Channel1: 2.0}, {Channel1: 3.0},
	}
	if err := emgRepo.AddSamplesBatch(ctx, session.ID, batchSamples); err != nil {
		t.Fatalf("AddSamplesBatch failed: %v", err)
	}

	count, err := emgRepo.CountSamplesBySessionID(ctx, session.ID)
	if err != nil {
		t.Fatalf("CountSamplesBySessionID failed: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 samples, got %d", count)
	}

	samples, err := emgRepo.FindSamplesBySessionID(ctx, session.ID, 1, 10)
	if err != nil {
		t.Fatalf("FindSamplesBySessionID failed: %v", err)
	}
	if len(samples) != 4 {
		t.Errorf("expected 4 samples, got %d", len(samples))
	}
}

func TestTrainingFileRepository(t *testing.T) {
	userRepo := NewUserRepository(testDB.Pool)
	deviceRepo := NewDeviceRepository(testDB.Pool)
	fileRepo := NewTrainingFileRepository(testDB.Pool)
	ctx := context.Background()

	userRepo.Create(ctx, "file-user", "f@test.com")
	device, _ := deviceRepo.Create(ctx, "file-user", &model.CreateDeviceRequest{
		Name: "Band", Type: "sensor", HWVersion: "1.0",
	})

	f, err := fileRepo.Create(ctx, "file-user", device.ID, "data.csv", "/tmp/data.csv", "test", 1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if f.OriginalName != "data.csv" {
		t.Errorf("expected data.csv, got %s", f.OriginalName)
	}

	found, err := fileRepo.FindByID(ctx, f.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.ID != f.ID {
		t.Errorf("expected %s, got %s", f.ID, found.ID)
	}

	if err := fileRepo.Delete(ctx, f.ID, "file-user"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = fileRepo.FindByID(ctx, f.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestStatsRepository(t *testing.T) {
	userRepo := NewUserRepository(testDB.Pool)
	deviceRepo := NewDeviceRepository(testDB.Pool)
	emgRepo := NewEMGRepository(testDB.Pool)
	trainingRepo := NewTrainingRepository(testDB.Pool)
	statsRepo := NewStatsRepository(testDB.Pool)
	ctx := context.Background()

	userRepo.Create(ctx, "stats-user", "st@test.com")
	device, _ := deviceRepo.Create(ctx, "stats-user", &model.CreateDeviceRequest{
		Name: "Band", Type: "prosthetic", HWVersion: "1.0",
	})
	sess, _ := emgRepo.CreateSession(ctx, "stats-user", &model.CreateEMGSessionRequest{
		DeviceID: device.ID, Label: "stats-sess",
	})
	job, err := trainingRepo.Create(ctx, "stats-user", &model.CreateTrainingJobRequest{
		SessionIDs: []string{sess.ID},
	})
	if err != nil {
		t.Fatalf("Create training job failed: %v", err)
	}
	trainingRepo.UpdateStatus(ctx, job.ID, model.TrainingStatusCompleted, "/model.h5", 0.9, "")

	stats, err := statsRepo.GetDashboardStats(ctx, "stats-user")
	if err != nil {
		t.Fatalf("GetDashboardStats failed: %v", err)
	}
	if stats.DeviceCount != 1 {
		t.Errorf("expected 1 device, got %d", stats.DeviceCount)
	}
	if stats.CompletedTrainings != 1 {
		t.Errorf("expected 1 completed training, got %d", stats.CompletedTrainings)
	}
	if stats.TotalTrainings != 1 {
		t.Errorf("expected 1 total training, got %d", stats.TotalTrainings)
	}
	if stats.AverageAccuracy != 0.9 {
		t.Errorf("expected 0.9 accuracy, got %f", stats.AverageAccuracy)
	}
}
