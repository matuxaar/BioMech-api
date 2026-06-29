import logging
import os

from dotenv import load_dotenv
from fastapi import FastAPI
from starlette.responses import JSONResponse
from app.router import train, predict
from app.config import settings

load_dotenv()

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger(__name__)
logger.info("starting ML service")

app = FastAPI(
    title="Desertacia ML Service",
    description="EMG signal processing and muscle contraction classification",
    version="0.1.0",
)

app.include_router(train.router)
app.include_router(predict.router)


@app.get("/health")
def health():
    db_ok = False
    conn = None
    try:
        import psycopg2
        conn = psycopg2.connect(settings.database_url, connect_timeout=3)
        conn.close()
        conn = None
        db_ok = True
    except Exception as e:
        logger.warning("health check: database connection failed: %s", e)
    finally:
        if conn is not None:
            try:
                conn.close()
            except Exception:
                pass
    if db_ok:
        return JSONResponse(content={"status": "ok", "database": "connected"}, status_code=200)
    return JSONResponse(content={"status": "degraded", "database": "unavailable"}, status_code=503)
