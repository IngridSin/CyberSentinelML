import threading

from flask import Flask, request, jsonify
import sys
import os
#sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from Common.db import DB, ssh_tunnel
from Emails.dkim import validate_email_headers
from Emails.imUtils import get_email_by_message_id, run_ml_model, validate_headers, update_email_analysis
from dotenv import load_dotenv

import atexit

from Redis.redisWorker import run_flow_worker

load_dotenv()
app = Flask(__name__)


atexit.register(lambda: ssh_tunnel and ssh_tunnel.stop())


@app.before_request
def log_request_info():
    print("Request Method:", request.method)
    print("Request Headers:", request.headers)
    print("Request Body:", request.get_data())


@app.route("/predict", methods=["POST"])
def predict():
    data = request.get_json()
    message_id = data.get("message_id")
    print(data, message_id)
    if not message_id:
        return jsonify({"error": "Missing message_id"}), 400

    #  Get email from DB
    email = get_email_by_message_id(db_im, message_id)
    print("message", email, message_id)
    if not email:
        return jsonify({"error": "Email not found"}), 404

    # Run ML model
    ml_result = run_ml_model(email)

    # Run Header Validation
    header_result = validate_email_headers(email)

    # Merge both results
    combined_result = {
        **ml_result,
        **header_result,
        "message_id": message_id
    }

    # Update database with everything
    update_email_analysis(db_im, combined_result)
    print("=====RESULTS=====", "\n",combined_result)


    return jsonify({
        "message_id": message_id,
    }), 200

@app.route("/health", methods=["GET"])
def health_check():
    return jsonify({
        "status": "ok",
        "message": "Flask server is running!"
    }), 200

@app.errorhandler(403)
def forbidden(e):
    return jsonify({"error": "Access forbidden", "details": str(e)}), 403


if __name__ == "__main__":
    db_pc = DB(os.getenv("DB_NAME"))
    db_im = DB(db_name=os.getenv("IM_DB_NAME"))

    try:
        print("Starting Flask ML API server...")
        db_pc.test_query()

        threading.Thread(target=run_flow_worker, args=(db_pc,), daemon=True).start()

        app.run(debug=True, use_reloader=False, host="0.0.0.0", port=5050)

    except KeyboardInterrupt:
        print("Shutting down Flask server...")
        db_im.close()
        db_pc.close()
