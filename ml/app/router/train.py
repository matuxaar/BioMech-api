from fastapi import APIRouter, HTTPException

from app.models.schemas import TrainRequest, TrainResponse
from app.models.trainer import train
from app.services.data_loader import fetch_samples

router = APIRouter()


@router.post("/train", response_model=TrainResponse)
def train_endpoint(req: TrainRequest):
    data = fetch_samples(req.session_ids)
    if len(data) == 0:
        raise HTTPException(
            status_code=400,
            detail=f"no samples found for sessions: {req.session_ids}",
        )

    try:
        model_path, accuracy = train(data, labels=None, job_id=req.job_id)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

    return TrainResponse(
        job_id=req.job_id,
        status="completed",
        accuracy=accuracy,
        model_path=model_path,
    )
