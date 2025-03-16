import json

from flask import Flask, jsonify, request
import pyshark
import threading

app = Flask(__name__)

capture_thread = None
captured_packets = []
stop_capture = False

def capture_packets(interface="eth0", max_packets=100):
    """Captures packets and stores them in memory."""
    global captured_packets, stop_capture
    captured_packets = []

    cap = pyshark.LiveCapture(interface=interface)

    for packet in cap.sniff_continuously(packet_count=max_packets):
        if stop_capture:
            break
        packet_info = {
            "timestamp": str(packet.sniff_time),
            "length": packet.length,
            "protocol": packet.highest_layer,
            "src": packet.ip.src if hasattr(packet, 'ip') else "N/A",
            "dst": packet.ip.dst if hasattr(packet, 'ip') else "N/A",
        }
        print("Packet:",packet)
        captured_packets.append(packet_info)

    cap.close()

@app.route('/start_capture', methods=['POST'])
def start_capture():
    """API endpoint to start packet capture."""
    global capture_thread, stop_capture

    if capture_thread and capture_thread.is_alive():
        return jsonify({"status": "error", "message": "Capture already running"}), 400

    data = request.json
    interface = data.get("interface", "eth0")
    max_packets = int(data.get("max_packets", 100))

    stop_capture = False
    capture_thread = threading.Thread(target=capture_packets, args=(interface, max_packets))
    capture_thread.start()

    return jsonify({"status": "success", "message": f"Started capturing on {interface}"})

@app.route('/stop_capture', methods=['POST'])
def stop_capture_fn():
    """API endpoint to stop packet capture."""
    global stop_capture

    stop_capture = True
    f = open("captured_packets.txt", "w")
    f.write(json.dumps(captured_packets))
    return jsonify({"status": "success", "message": "Stopping capture..."})

@app.route('/packets', methods=['GET'])
def get_packets():
    """API endpoint to get captured packets."""
    return jsonify({"status": "success", "packets": captured_packets})

@app.route('/')
def home():
    return jsonify({"message": "Pyshark Flask API is running!"})

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=9091, debug=True)
