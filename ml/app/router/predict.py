import logging

from fastapi import APIRouter, HTTPException

from app.models.schemas import (
    PredictRequest,
    PredictResponse,
    StreamPredictRequest,
    StreamPredictResponse,
)
from app.services.predictor import Predictor
from app.config import settings

logger = logging.getLogger(__name__)

router = APIRouter()
predictor = Predictor()


def _predict_samples(samples):
    if not samples:
        raise HTTPException(status_code=400, detail="no samples provided")
    try:
        predictions = predictor.predict(samples)
    except RuntimeError as e:
        raise HTTPException(status_code=503, detail=str(e))
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.exception("predict failed")
        raise HTTPException(status_code=500, detail=str(e))
    return predictions


@router.post("/predict", response_model=PredictResponse)
def predict_endpoint(req: PredictRequest):
    predictions = _predict_samples(req.samples)
    return PredictResponse(predictions=predictions)


@router.post("/predict/stream", response_model=StreamPredictResponse)
def predict_stream_endpoint(req: StreamPredictRequest):
    predictions = _predict_samples(req.samples)
    return StreamPredictResponse(predictions=predictions)
