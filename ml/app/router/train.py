import logging

import httpx
from fastapi import APIRouter, HTTPException
from tenacity import retry, stop_after_attempt, wait_exponential, retry_if_exception_type

from app.models.schemas import TrainRequest, TrainResponse
from app.models.trainer import train
from app.services.data_loader import fetch_samples_with_labels
from app.config import settings

logger = logging.getLogger(__name__)


@retry(
    stop=stop_after_attempt(3),
    wait=wait_exponential(multiplier=1, min=1, max=10),
    retry=retry_if_exception_type((httpx.HTTPError, httpx.TimeoutException)),
    before_sleep=lambda rs: logger.warning(
        "callback retry %d", rs.attempt_number
    ),
)
def _notify_backend(job_id: str, status: str, accuracy: float | None = None, error_message: str | None = None):
    url = f"{settings.backend_callback_url}/{job_id}/status"
    headers = {"X-API-Key": settings.backend_api_key}
    payload = {"job_id": job_id, "status": status}
    if accuracy is not None:
        payload["accuracy"] = accuracy
    if error_message is not None:
        payload["error_message"] = error_message
    resp = httpx.post(url, json=payload, headers=headers, timeout=10)
    resp.raise_for_status()
    logger.info("callback for job %s: status=%s response=%d", job_id, status, resp.status_code)


def _try_notify(*args, **kwargs):
    try:
        _notify_backend(*args, **kwargs)
    except Exception as e:
        logger.error("callback failed after retries: %s", e)

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
        _try_notify(req.job_id, "completed", accuracy=accuracy)
    except Exception as e:
        _try_notify(req.job_id, "failed", error_message=str(e))
        raise HTTPException(status_code=500, detail=str(e))

    return TrainResponse(
        job_id=req.job_id,
        status="completed",
        accuracy=accuracy,
        model_path=model_path,
    )
