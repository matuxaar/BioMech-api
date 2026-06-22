from fastapi import APIRouter, HTTPException

from app.models.schemas import PredictRequest, PredictResponse
from app.services.predictor import Predictor

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
