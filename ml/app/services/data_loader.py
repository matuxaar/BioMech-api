import logging
from contextlib import contextmanager

import psycopg2
from psycopg2.pool import ThreadedConnectionPool
import numpy as np
from app.config import settings, GESTURE_LABELS

_pool = None


def _get_pool():
    global _pool
    if _pool is None:
        _pool = ThreadedConnectionPool(1, 10, settings.database_url)
    return _pool


@contextmanager
def _get_conn():
    pool = _get_pool()
    conn = pool.getconn()
    try:
        yield conn
    finally:
        pool.putconn(conn)

logger = logging.getLogger(__name__)


def _label_to_index(label: str) -> int:
    clean = label.strip().lower()
    if clean in GESTURE_LABELS:
        return GESTURE_LABELS.index(clean)
    logger.warning("unknown label '%s', mapping to 0 (rest)", label)
    return 0


def fetch_samples_with_labels(session_ids: list[str]) -> tuple[np.ndarray, np.ndarray]:
    if not session_ids:
        logger.warning("fetch_samples_with_labels called with empty session_ids")
        return np.empty((0, settings.n_channels)), np.empty((0,), dtype=np.int32)

    placeholders = ",".join(["%s"] * len(session_ids))

    query = f"""
        SELECT s.session_id, s.channel_1, s.channel_2, s.channel_3, s.channel_4,
               s.channel_5, s.channel_6, s.channel_7, s.channel_8,
               COALESCE(e.label, '')
        FROM emg_samples s
        JOIN emg_sessions e ON s.session_id = e.id
        WHERE s.session_id IN ({placeholders})
        ORDER BY s.session_id, s.timestamp ASC
    """

    try:
        with _get_conn() as conn:
            with conn.cursor() as cur:
                cur.execute(query, session_ids)
                rows = cur.fetchall()
    except Exception as e:
        logger.error("database query failed: %s", e)
        raise

    if not rows:
        logger.info("no samples found for session_ids: %s", session_ids)
        return np.empty((0, settings.n_channels)), np.empty((0,), dtype=np.int32)

    samples = np.array([r[1:9] for r in rows], dtype=np.float32)
    labels = np.array(
        [_label_to_index(r[9]) for r in rows],
        dtype=np.int32,
    )

    logger.info("loaded %d samples from %d sessions", len(samples), len(session_ids))
    return samples, labels
