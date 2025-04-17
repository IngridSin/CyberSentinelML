import psycopg2
import select
import threading
import time
import os
from psycopg2.extras import execute_values
from sshtunnel import SSHTunnelForwarder
from dotenv import load_dotenv

load_dotenv()

# Shared global SSH tunnel
ssh_tunnel = None
tunnel_lock = threading.Lock()

def get_ssh_tunnel():
    global ssh_tunnel
    with tunnel_lock:
        if ssh_tunnel is None or not ssh_tunnel.is_active:
            print("Opening shared SSH tunnel...")
            ssh_tunnel = SSHTunnelForwarder(
                (os.getenv("SSH_HOST"), int(os.getenv("SSH_PORT"))),
                ssh_username=os.getenv("SSH_USER"),
                ssh_password=os.getenv("SSH_PASSWORD"),
                ssh_private_key=os.getenv("SSH_PRIVATE_KEY"),
                remote_bind_address=(os.getenv("DB_HOST"), int(os.getenv("DB_PORT")))
            )
            ssh_tunnel.start()
            print("SSH Tunnel established on local port:", ssh_tunnel.local_bind_port)
        return ssh_tunnel


class DB:
    def __init__(self, db_name):
        self.db_name = db_name
        self.conn = None
        self.cursor = None
        self.flow_id_queue = set()
        self.queue_lock = threading.Lock()
        self.BATCH_INTERVAL = 10

        self._connect_db()

    def _connect_db(self):
        try:
            tunnel = get_ssh_tunnel()

            self.conn = psycopg2.connect(
                dbname=self.db_name,
                user=os.getenv("DB_USER"),
                password=os.getenv("DB_PASSWORD"),
                host="127.0.0.1",
                port=tunnel.local_bind_port
            )
            self.conn.set_isolation_level(psycopg2.extensions.ISOLATION_LEVEL_AUTOCOMMIT)
            self.cursor = self.conn.cursor()
            print(f"Connected to {self.db_name} through shared SSH tunnel!")
        except Exception as e:
            print(f"Failed to connect to DB '{self.db_name}': {e}")
            self.close()

    def test_query(self):
        try:
            self.cursor.execute("SELECT version();")
            version = self.cursor.fetchone()
            print(f"Connected to PostgreSQL. Version: {version[0]}")
        except Exception as e:
            print(f"Test query failed: {e}")

    def close(self):
        if self.cursor:
            self.cursor.close()
        if self.conn:
            self.conn.close()
        print(f"Closed DB connection to {self.db_name}")



    def insert_predicted_batch(self, df):
        raw_columns = [
            'Flow ID', 'Source IP', 'Destination IP', 'Source Port', 'Destination Port', 'Protocol',
            'Flow Duration', 'Total Fwd Packets', 'Total Backward Packets',
            'Total Length of Fwd Packets', 'Total Length of Bwd Packets',
            'Fwd Packet Length Max', 'Fwd Packet Length Min', 'Fwd Packet Length Mean', 'Fwd Packet Length Std',
            'Bwd Packet Length Max', 'Bwd Packet Length Min', 'Bwd Packet Length Mean', 'Bwd Packet Length Std',
            'Flow Bytes/s', 'Flow Packets/s',
            'Flow IAT Mean', 'Flow IAT Std', 'Flow IAT Max', 'Flow IAT Min',
            'Fwd IAT Total', 'Fwd IAT Mean', 'Fwd IAT Std', 'Fwd IAT Max', 'Fwd IAT Min',
            'Bwd IAT Total', 'Bwd IAT Mean', 'Bwd IAT Std', 'Bwd IAT Max', 'Bwd IAT Min',
            'Fwd PSH Flags', 'Bwd PSH Flags', 'Fwd URG Flags', 'Bwd URG Flags',
            'Fwd Header Length', 'Bwd Header Length', 'Fwd Packets/s', 'Bwd Packets/s',
            'Min Packet Length', 'Max Packet Length', 'Packet Length Mean', 'Packet Length Std', 'Packet Length Variance',
            'FIN Flag Count', 'SYN Flag Count', 'RST Flag Count', 'PSH Flag Count', 'ACK Flag Count', 'URG Flag Count',
            'CWE Flag Count', 'ECE Flag Count',
            'Down/Up Ratio', 'Average Packet Size', 'Avg Fwd Segment Size', 'Avg Bwd Segment Size',
            'Fwd Avg Bytes/Bulk', 'Fwd Avg Packets/Bulk', 'Fwd Avg Bulk Rate',
            'Bwd Avg Bytes/Bulk', 'Bwd Avg Packets/Bulk', 'Bwd Avg Bulk Rate',
            'Subflow Fwd Packets', 'Subflow Fwd Bytes', 'Subflow Bwd Packets', 'Subflow Bwd Bytes',
            'act_data_pkt_fwd', 'min_seg_size_forward',
            'Active Mean', 'Active Std', 'Active Max', 'Active Min',
            'Idle Mean', 'Idle Std', 'Idle Max', 'Idle Min',
            'XGBoost Prediction'

        ]

        # Normalize: match DB field names
        def normalize_column(col):
            return col.replace("/", "_per_").replace(" ", "_").replace(".", "").replace("-", "_").lower()


        db_columns = [normalize_column(col) for col in raw_columns]


        df.columns = [normalize_column(c) for c in df.columns]

        rows = df[db_columns].values.tolist()

        with self.conn.cursor() as cur:
            execute_values(
                cur,
                f"""
                    INSERT INTO test_schema.test_network_flow ({', '.join(db_columns)})
                    VALUES %s
                    """,
                rows
            )
        self.conn.commit()
