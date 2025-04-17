import redis
import json
import time
from datetime import datetime
import traceback


import pandas as pd
import os

from Network.featureExtractor import extract_cic_features
from Network.predictPackets import predict_packets


BATCH_SIZE = 50
r = redis.Redis(host=os.getenv("REDIS_ADDR"), port=6379, db=0)

def parse_timestamps(ts):
    return [datetime.fromisoformat(t) if isinstance(t, str) else t for t in ts]


def packets_processing(db, batch, flows):
    cic_df = pd.DataFrame([extract_cic_features(flow) for flow in flows])
    result_df = predict_packets(cic_df)

    print(result_df[["Flow ID", "XGBoost Prediction"]].head())

    db.insert_predicted_batch(result_df)

    batch.clear()
    flows.clear()


def run_flow_worker(db):
    try:
        print("Starting Flow Worker in background...")
        batch = []
        flows = []

        while True:
            data = r.lpop("packet_queue")
            if data:
                try:
                    pkt = json.loads(data)
                    flows.append(pkt)
                    batch.append(pkt)

                    # print("\nRetrieved Packet:", data, "\n")

                    if len(batch) >= BATCH_SIZE:
                        # print(f"Processing batch of {len(batch)} flows...")
                        packets_processing(db, batch, flows)


                except Exception as e:
                    print("Error processing flow packet:", e)
                    traceback.print_exc()
            else:
                if batch:
                    print("Flushing timed-out batch...")
                    packets_processing(db, batch, flows)
                    time.sleep(0.1)

    except Exception as e:
        print("Worker failed to start:", e)
        traceback.print_exc()
