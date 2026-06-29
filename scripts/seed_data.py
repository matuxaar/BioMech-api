"""
Test EMG data generator.

Creates a user, devices, EMG sessions and synthetic samples
for all 5 gesture classes: rest, fist, open, pinch, point.

Usage (run with DEV_MODE=true on the backend):
    pip install requests numpy
    python scripts/seed_data.py [--api-url http://localhost:8080]
"""

import argparse
import math
import random
from datetime import datetime, timezone

import requests

API_URL = "http://localhost:8080"
USER_EMAIL = "test@desertacia.dev"

GESTURES = ["rest", "fist", "open", "pinch", "point"]
N_SAMPLES_PER_SESSION = 500
N_CHANNELS = 8
SAMPLING_RATE = 100  # Hz


def generate_emg_samples(gesture: str, n: int) -> list[dict]:
    samples = []
    t0 = datetime.now(timezone.utc).timestamp()

    # Signal parameters per gesture
    params = {
        "rest":  {"amp": 0.01, "freq": 5,   "noise": 0.005},
        "fist":  {"amp": 0.8,  "freq": 50,  "noise": 0.05},
        "open":  {"amp": 0.5,  "freq": 30,  "noise": 0.04},
        "pinch": {"amp": 0.6,  "freq": 60,  "noise": 0.03},
        "point": {"amp": 0.4,  "freq": 40,  "noise": 0.04},
    }
    p = params[gesture]

    for i in range(n):
        t = i / SAMPLING_RATE
        channels = []
        for ch in range(N_CHANNELS):
            phase = random.uniform(0, 2 * math.pi)
            value = (
                p["amp"] * math.sin(2 * math.pi * p["freq"] * t + phase)
                + p["amp"] * 0.3 * math.sin(2 * math.pi * p["freq"] * 2 * t + phase)
                + random.gauss(0, p["noise"])
            )
            channels.append(round(value, 6))

        ts = datetime.fromtimestamp(t0 + i / SAMPLING_RATE, tz=timezone.utc).isoformat()
        samples.append({
            "timestamp": ts,
            "channel_1": channels[0],
            "channel_2": channels[1],
            "channel_3": channels[2],
            "channel_4": channels[3],
            "channel_5": channels[4],
            "channel_6": channels[5],
            "channel_7": channels[6],
            "channel_8": channels[7],
        })

    return samples


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--api-url", default=API_URL)
    args = parser.parse_args()
    base = args.api_url

    session = requests.Session()
    DUMMY_TOKEN = "dev-mode-token"

    # --- 1. Sync user (Firebase sync, in DEV_MODE it creates the user) ---
    print("[1/6] Syncing user...")
    r = session.post(f"{base}/api/v1/auth/firebase", headers={"Authorization": f"Bearer {DUMMY_TOKEN}"})
    if r.status_code == 200:
        print("      User created")
        user = r.json()
        print(f"      User: {user.get('id', 'unknown')}")
    else:
        print(f"      Failed to sync user: {r.status_code} {r.text}")
        return

    session.headers.update({"Authorization": f"Bearer {DUMMY_TOKEN}"})

    # --- 2. Create sensor device ---
    print("[2/6] Creating sensor device...")
    r = session.post(f"{base}/api/v1/devices", json={
        "type": "sensor",
        "name": "EMG Armband v1",
        "hw_version": "1.0.0",
    })
    if r.status_code == 201:
        sensor_id = r.json()["id"]
        print(f"      Sensor: {sensor_id}")
    else:
        print(f"      Error: {r.status_code} {r.text}")
        return

    # --- 3. Create prosthetic device ---
    print("[3/6] Creating prosthetic device...")
    r = session.post(f"{base}/api/v1/devices", json={
        "type": "prosthetic",
        "name": "Bionic Hand Pro",
        "hw_version": "2.0.0",
    })
    if r.status_code == 201:
        prosthetic_id = r.json()["id"]
        print(f"      Prosthetic: {prosthetic_id}")
    else:
        print(f"      Error: {r.status_code} {r.text}")
        return

    # --- 4-5. For each gesture: create session + upload samples ---
    print("[4/6] Creating EMG sessions...")
    session_ids = {}
    for gesture in GESTURES:
        r = session.post(f"{base}/api/v1/emg/sessions", json={
            "device_id": sensor_id,
            "label": gesture,
        })
        if r.status_code == 201:
            sid = r.json()["id"]
            session_ids[gesture] = sid
            print(f"      {gesture}: {sid}")
        else:
            print(f"      Error: {r.status_code} {r.text}")
            return

    print("[5/6] Generating and uploading samples...")
    for gesture in GESTURES:
        sid = session_ids[gesture]
        samples = generate_emg_samples(gesture, N_SAMPLES_PER_SESSION)

        r = session.post(f"{base}/api/v1/emg/sessions/{sid}/samples/batch", json={
            "samples": samples,
        })
        if r.status_code == 201:
            print(f"      {gesture}: {len(samples)} samples uploaded")
        else:
            print(f"      Error: {r.status_code} {r.text}")
            return

    # --- 6. End sessions ---
    print("[6/6] Ending sessions...")
    for gesture in GESTURES:
        sid = session_ids[gesture]
        r = session.post(f"{base}/api/v1/emg/sessions/{sid}/end")
        if r.status_code == 200:
            print(f"      {gesture}: ended")
        else:
            print(f"      Error: {r.status_code} {r.text}")

    print("\nDone! Data loaded:")
    print(f"  User: {USER_EMAIL}")
    print(f"  Sessions: {len(session_ids)}")
    print(f"  Samples: {len(GESTURES) * N_SAMPLES_PER_SESSION}")
    print(f"\nNow you can start training:")
    print(f"  curl -X POST {base}/api/v1/training/jobs \\")
    print(f"    -H 'Authorization: Bearer {DUMMY_TOKEN}' \\")
    print(f"    -H 'Content-Type: application/json' \\")
    print(f"    -d '{{\"session_ids\": [\"{session_ids['fist']}\", \"{session_ids['open']}\"]}}'")


if __name__ == "__main__":
    main()
