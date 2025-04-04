from flask import Flask, request, jsonify
from FlaskServer.Common.db import DB
from FlaskServer.Emails.dkim import validate_email_headers
from FlaskServer.Emails.imUtils import get_email_by_message_id, run_ml_model, validate_headers, update_email_analysis
from dotenv import load_dotenv
import os
import atexit

load_dotenv()
app = Flask(__name__)

# Initialize DB
db = DB(db_name=os.getenv("IM_DB_NAME"))

# Register graceful shutdown
atexit.register(db.close)

@app.before_request
def log_request_info():
    print("Request Method:", request.method)
    print("Request Headers:", request.headers)
    print("Request Body:", request.get_data())


@app.route("/predict", methods=["POST"])
def predict():
    print("PREDICTING:")
    data = request.get_json()
    message_id = data.get("message_id")
    print(data, message_id)
    if not message_id:
        return jsonify({"error": "Missing message_id"}), 400

    #  Get email from DB
    email = get_email_by_message_id(db, message_id)
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
    update_email_analysis(db, combined_result)
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
    try:
        print("Starting Flask ML API server...")
        app.run(debug=True, host="0.0.0.0", port=5000)
    except KeyboardInterrupt:
        print("Shutting down Flask server...")
        db.close()
