import psycopg2
import numpy as np
from app.config import settings, GESTURE_LABELS


def _label_to_index(label: str) -> int:
    clean = label.strip().lower()
    if clean in GESTURE_LABELS:
        return GESTURE_LABELS.index(clean)
    return 0


def fetch_samples_with_labels(session_ids: list[str]) -> tuple[np.ndarray, np.ndarray]:
    if not session_ids:
        return np.empty((0, settings.n_channels)), np.empty((0,), dtype=np.int32)

    placeholders = ",".join(["%s"] * len(session_ids))

    with psycopg2.connect(settings.database_url) as conn:
        with conn.cursor() as cur:
            cur.execute(
                f"SELECT id, COALESCE(label, '') FROM emg_sessions WHERE id IN ({placeholders})",
                session_ids,
            )
            session_map = {row[0]: row[1] for row in cur.fetchall()}

            query = f"""
                SELECT session_id, channel_1, channel_2, channel_3, channel_4,
                       channel_5, channel_6, channel_7, channel_8
                FROM emg_samples
                WHERE session_id IN ({placeholders})
                ORDER BY timestamp ASC
            """
            cur.execute(query, session_ids)
            rows = cur.fetchall()

    if not rows:
        return np.empty((0, settings.n_channels)), np.empty((0,), dtype=np.int32)

    samples = np.array([r[1:] for r in rows], dtype=np.float32)
    labels = np.array(
        [_label_to_index(session_map.get(r[0], "")) for r in rows],
        dtype=np.int32,
    )

    return samples, labels
