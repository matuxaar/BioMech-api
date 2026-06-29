# Desertacia

Bionic prosthetic control system based on EMG signal processing. The system reads muscle electrical activity from an armband sensor (8 channels), classifies gestures via a TensorFlow model, and sends control commands to a prosthetic hand in real time.

## Architecture

```
┌─────────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  Android App        │────▶│  Go Backend       │────▶│  ML Service       │
│  (Kotlin/MVI/Compose)│     │  (Gin/Pgx)        │     │  (FastAPI/TF)     │
│                     │     │                   │     │                  │
│  - BLE armband      │     │  - REST API       │     │  - Gesture       │
│  - EMG streaming    │     │  - WebSocket      │     │    prediction    │
│  - Dashboard        │     │  - Auth (Firebase)│     │  - Model         │
│  - Training UI      │     │  - File uploads   │     │    training      │
│  - Offline queue    │     │  - PostgreSQL     │     │  - GCS storage   │
└─────────────────────┘     └──────────────────┘     └──────────────────┘
                                    │
                                    ▼
                            ┌──────────────────┐
                            │  PostgreSQL       │
                            │                   │
                            │  - Users         │
                            │  - Devices       │
                            │  - EMG sessions  │
                            │  - Training jobs │
                            └──────────────────┘
```

## How it works

1. **EMG Data Collection** — The app connects to an 8-channel armband sensor via BLE (or simulated). EMG samples are streamed to the backend in real time via WebSocket.

2. **Gesture Classification** — The ML service runs a TensorFlow model that classifies incoming EMG samples into 5 gesture classes:
   - `rest` — no muscle activation
   - `fist` — hand clenched
   - `open` — hand open
   - `pinch` — thumb + index pinch
   - `point` — index finger extended

3. **Model Training** — Users can record EMG sessions for each gesture, then trigger training. The backend sends samples to the ML service, which trains a neural network and saves the model. The ML service notifies the backend when training completes.

4. **Prosthetic Control** — Classified gestures are sent as commands to a prosthetic hand device via BLE. Commands (open/close/pinch/point) are configured per device.

5. **Offline Support** — The app queues failed requests locally (Room) and replays them when connectivity is restored.

## Project Structure

```
desertacia/
├── backend/              # Go backend (Gin framework)
│   ├── cmd/api/          # Entry point
│   ├── internal/
│   │   ├── config/       # Environment config
│   │   ├── handler/      # HTTP handlers
│   │   ├── middleware/   # Auth, CORS
│   │   ├── migrations/   # SQL migrations
│   │   ├── model/        # Domain models
│   │   ├── repository/   # Data access (pgx)
│   │   └── service/      # Business logic
│   ├── migrations/       # SQL migration files
│   └── Dockerfile
├── ml/                   # Python ML service (FastAPI)
│   ├── app/
│   │   ├── models/       # TF model definition, training
│   │   ├── router/       # API routes
│   │   └── services/     # Data loading, feature extraction, prediction
│   ├── Dockerfile
│   └── requirements.txt
├── BioMech-app/          # Android/iOS app (KMP)
│   ├── androidApp/       # Android entry point
│   ├── iosApp/           # iOS entry point
│   ├── shared/           # Shared Compose Multiplatform UI
│   ├── core/             # Core modules (network, database, DI, BLE)
│   └── feature/          # Feature modules (auth, devices, training, dashboard)
├── scripts/
│   └── seed_data.py      # Test data generation script
├── docker-compose.yml    # Root compose (all services)
└── cloudbuild.yaml       # Google Cloud Build CI/CD
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.25+ (for local backend development)
- Python 3.12+ (for local ML development)
- Android Studio (for app development)
- Firebase project (for auth)

### Run all services

```bash
# Start PostgreSQL, API, and ML service
docker compose up --build

# The API is available at http://localhost:8080
# The ML service is available at http://localhost:8000
```

### Seed test data

```bash
pip install requests numpy
python scripts/seed_data.py
```

This creates a user, sensor/prosthetic devices, and 500 EMG samples for each of the 5 gesture classes.

### Environment Variables

Key environment variables (see `backend/.env.example` and `ml/.env.example`):

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | API server port |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/desertacia?sslmode=disable` | PostgreSQL connection |
| `ML_SERVICE_URL` | `http://ml:8000` | ML service URL |
| `INTERNAL_API_KEY` | (required) | Shared key for ML→API callbacks |
| `DEV_MODE` | `true` | Skip Firebase auth, use dev stub |
| `CORS_ORIGINS` | `http://localhost:8080` | Comma-separated allowed CORS origins |
| `UPLOADS_DIR` | `./uploads` | Upload storage directory |
| `MAX_UPLOAD_SIZE_MB` | `50` | Max upload file size in MB |

### Backend (local)

```bash
cd backend
go mod download
go run ./cmd/api
```

### ML Service (local)

```bash
cd ml
pip install -r requirements.txt
uvicorn app.main:app --reload
```

## API Overview

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | No | Health check |
| `POST` | `/api/v1/auth/firebase` | Firebase | Sync Firebase user |
| `GET` | `/api/v1/me` | JWT | Get current user |
| `PUT` | `/api/v1/me` | JWT | Update profile |
| `POST` | `/api/v1/me/avatar` | JWT | Upload avatar |
| `GET/POST/PUT/DELETE` | `/api/v1/devices` | JWT | Device CRUD |
| `POST` | `/api/v1/emg/sessions` | JWT | Start EMG session |
| `POST` | `/api/v1/emg/sessions/:id/samples/batch` | JWT | Upload EMG samples |
| `POST` | `/api/v1/training/jobs` | JWT | Start training |
| `GET` | `/api/v1/training/jobs` | JWT | List training jobs |
| `GET` | `/predict/ws` | JWT | WebSocket prediction stream |
| `POST` | `/api/v1/training/jobs/:id/status` | API key | ML → Backend callback |

## Deployment

The project uses Google Cloud Build for CI/CD. GitHub Actions workflows are also configured:

- **Backend CI**: lint, test (with PostgreSQL), build
- **ML CI**: lint (ruff), test (pytest), build
- **App CI**: Android debug build
- **Deploy**: Triggered on push to `main`, deploys to Cloud Run

For production:

```bash
# Set DEV_MODE=false, configure Firebase credentials,
# set CORS_ORIGINS to your domain, and provide a secure INTERNAL_API_KEY.
```

## License

Proprietary — all rights reserved.
