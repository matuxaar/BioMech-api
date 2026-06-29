from fastapi import APIRouter, HTTPException

from app.models.schemas import (
    PredictRequest,
    PredictResponse,
    StreamPredictRequest,
    StreamPredictResponse,
)
from app.services.predictor import Predictor
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

    predictions = predictor.predict(req.samples)

    return StreamPredictResponse(predictions=predictions)
