import psycopg2
import select
import threading
import time
import os
from sshtunnel import SSHTunnelForwarder
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class DB:
    def __init__(self, db_name):
        self.db_name = db_name
        self.ssh_tunnel = None
        self.conn = None
        self.cursor = None
        self.flow_id_queue = set()
        self.queue_lock = threading.Lock()
        self.BATCH_INTERVAL = 10  # Process every 10 seconds

        self._connect_db()
        self._listen_for_notifications()
        # self._start_batch_processing_thread()

    def _open_ssh_tunnel(self):
        connection_success = False

        print("Starting SSH Tunnel...")

        while not connection_success:
            try:
                self.ssh_tunnel = SSHTunnelForwarder(
                    (os.getenv("SSH_HOST"), int(os.getenv("SSH_PORT"))),
                    ssh_username=os.getenv("SSH_USER"),
                    ssh_password=os.getenv("SSH_PASSWORD"),
                    ssh_private_key=os.getenv("SSH_PRIVATE_KEY"),
                    remote_bind_address=(os.getenv("DB_HOST"), int(os.getenv("DB_PORT")))
                )
                connection_success = True

            except Exception as e:
                time.sleep(0.5)

        self.ssh_tunnel.start()
        print("SSH Tunnel established! Connecting to PostgreSQL...")


    def _connect_db(self):
        """Establish an SSH tunnel and connect to PostgreSQL."""
        try:
            self._open_ssh_tunnel()

            self.conn = psycopg2.connect(
                dbname=self.db_name,
                user=os.getenv("DB_USER"),
                password=os.getenv("DB_PASSWORD"),
                host="127.0.0.1",
                port=self.ssh_tunnel.local_bind_port
            )

            self.conn.set_isolation_level(psycopg2.extensions.ISOLATION_LEVEL_AUTOCOMMIT)
            self.cursor = self.conn.cursor()
            print("Connected to PostgreSQL through SSH tunnel!")

        except Exception as e:
            print(f"Failed to establish SSH tunnel or PostgreSQL connection: {e}")
            self.close()


    def test_query(self):
        """Runs a test query to verify the database connection."""
        try:
            self.cursor.execute("SELECT version();")
            version = self.cursor.fetchone()
            print(f"Connected to PostgreSQL. Version: {version[0]}")

        except Exception as e:
            print(f"Test query failed: {e}")


    def _listen_for_notifications(self):

        """Start listening for PostgreSQL notifications."""

        try:
            self.cursor.execute("LISTEN NEW_PACKET;")

            print("Listening for Notification NEW_PACKET from PostgreSQL...")
        except Exception as e:
            print(f"Error setting up LISTEN: {e}")


    def _start_batch_processing_thread(self):
        """Start a background thread to process packets in batches."""
        threading.Thread(target=self._fetch_batch_and_process, daemon=True).start()


    def _fetch_batch_and_process(self):
        """Fetch and process packets in batch every BATCH_INTERVAL seconds."""
        while True:
            time.sleep(self.BATCH_INTERVAL)

            with self.queue_lock:
                if not self.flow_id_queue:
                    continue  # No packets to process

                batch_flow_ids = list(self.flow_id_queue)
                self.flow_id_queue.clear()

            print(f"Processing batch of {len(batch_flow_ids)} packets")

            self.cursor.execute("""
                SELECT flow_id, source_ip, destination_ip, protocol 
                FROM your_schema.your_table WHERE flow_id = ANY(%s);
            """, (batch_flow_ids,))
            packets = self.cursor.fetchall()

            # processed_results = {row[0]: f"{row[1]}-{row[2]}-{row[3]}-processed" for row in packets}

            # Bulk update results
            # self._update_batch_in_db(processed_results)


    def _update_batch_in_db(self, results):
        """Update PostgreSQL with processed batch data."""
        cursor = self.conn.cursor()
        update_query = "UPDATE your_schema.your_table SET processed_data = %s WHERE flow_id = %s;"

        for flow_id, processed_data in results.items():
            cursor.execute(update_query, (processed_data, flow_id))

        self.conn.commit()
        print(f"Updated {len(results)} rows in PostgreSQL.")


    def listen_for_new_packets(self):
        """Continuously listen for new packets from PostgreSQL."""


        while True:

            print("Listening for new packets...")

            if select.select([self.conn], [], [], 5) == ([], [], []):
                continue  # Timeout, loop again

            self.conn.poll()
            while self.conn.notifies:
                notification = self.conn.notifies.pop(0)

                print(f"Received new packet from PostgreSQL: {notification}")


                # flow_id = notification.payload
                #
                # with self.queue_lock:
                #     self.flow_id_queue.add(flow_id) # Flow ID to be queried

                # print(f"Queued flow ID for batch processing: {flow_id}")


    def close(self):
        """Close database connection and SSH tunnel."""
        if self.cursor:
            self.cursor.close()
        if self.conn:
            self.conn.close()
        if self.ssh_tunnel:
            self.ssh_tunnel.stop()
        print("Database and SSH tunnel closed.")


