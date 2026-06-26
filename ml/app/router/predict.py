from fastapi import APIRouter, HTTPException

import numpy as np

from app.models.schemas import (
    PredictRequest,
    PredictResponse,
    StreamPredictRequest,
    StreamPredictResponse,
)
from app.services.predictor import Predictor, GESTURE_LABELS
from app.services.features import extract_features
from app.config import settings

router = APIRouter()
predictor = Predictor()


@router.post("/predict", response_model=PredictResponse)
def predict_endpoint(req: PredictRequest):
    if not req.samples:
        raise HTTPException(status_code=400, detail="no samples provided")

    try:
        predictions = predictor.predict(req.samples)
    except RuntimeError as e:
        raise HTTPException(status_code=503, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

    return PredictResponse(predictions=predictions)


@router.post("/predict/stream", response_model=StreamPredictResponse)
def predict_stream_endpoint(req: StreamPredictRequest):
    if not req.samples:
        raise HTTPException(status_code=400, detail="no samples provided")

    model = predictor.get_model()

    data = np.array(
        [[s.ch1, s.ch2, s.ch3, s.ch4, s.ch5, s.ch6, s.ch7, s.ch8] for s in req.samples],
        dtype=np.float32,
    )

    if len(data) >= settings.window_size:
        features = extract_features(data[-settings.window_size:]).reshape(1, -1)
    elif len(data) >= 1:
        padded = np.zeros((settings.window_size, settings.n_channels), dtype=np.float32)
        padded[-len(data):] = data
        features = extract_features(padded).reshape(1, -1)
    else:
        raise HTTPException(status_code=400, detail="at least one sample required")

    predictions = model.predict(features, verbose=0)
    predicted_indices = np.argmax(predictions, axis=1)

    labels = [
        GESTURE_LABELS[idx] if idx < len(GESTURE_LABELS) else "unknown"
        for idx in predicted_indices
    ]

    return StreamPredictResponse(predictions=labels)
