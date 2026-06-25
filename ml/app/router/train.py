from fastapi import APIRouter, HTTPException

from app.models.schemas import TrainRequest, TrainResponse
from app.models.trainer import train
from app.services.data_loader import fetch_samples_with_labels

router = APIRouter()


@router.post("/train", response_model=TrainResponse)
def train_endpoint(req: TrainRequest):
    data, labels = fetch_samples_with_labels(req.session_ids)
    if len(data) == 0:
        raise HTTPException(
            status_code=400,
            detail=f"no samples found for sessions: {req.session_ids}",
        )

    unique_labels = set(labels.tolist())
    n_classes = len(unique_labels) if len(unique_labels) > 1 else 5

    try:
        model_path, accuracy = train(data, labels=labels, job_id=req.job_id, n_classes=n_classes)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

    return TrainResponse(
        job_id=req.job_id,
        status="completed",
        accuracy=accuracy,
        model_path=model_path,
    )
