import json

from FlaskServer.Common.db import DB
from dotenv import load_dotenv
import os
from typing import Dict, Any
import atexit

from FlaskServer.Emails.dkim import validate_email_headers
from FlaskServer.Emails.ml import predict_phishing

load_dotenv()


db = DB(os.getenv("IM_DB_NAME"))


def get_email_by_message_id(db, message_id):
    query = "SELECT * FROM emails WHERE message_id = %s"
    db.cursor.execute(query, (message_id,))
    row = db.cursor.fetchone()

    if row:
        return {
            "id": row[0],
            "message_id": row[1],
            "date": row[2],
            "subject": row[3],
            "sender": row[4],
            "recipient": row[5],
            "received": row[6],
            "return_path": row[7],
            "delivered_to": row[8],
            "body": row[9],
            "dkim": row[10],
            "spf": row[11],
            "attachments": row[12]
        }

    return None


def update_email_analysis(db, result):
    query = """
        UPDATE emails SET
            prediction = %s,
            winner_model = %s,
            winner_probability = %s,

            ensemble_predicted_label = %s,
            ensemble_prob_class0 = %s,
            ensemble_prob_class1 = %s,

            nb_pred = %s,
            rf_pred = %s,
            xgb_pred = %s,
            knn_pred = %s,
            logreg_pred = %s,

            nb_prob_class0 = %s,
            rf_prob_class0 = %s,
            xgb_prob_class0 = %s,
            knn_prob_class0 = %s,
            logreg_prob_class0 = %s,

            nb_prob_class1 = %s,
            rf_prob_class1 = %s,
            xgb_prob_class1 = %s,
            knn_prob_class1 = %s,
            logreg_prob_class1 = %s,

            top_5_words_from_nb = %s,
            top_5_words_from_rf = %s,
            top_5_words_from_xgb = %s,
            top_5_words_from_knn = %s,
            top_5_words_from_logreg = %s,

            risk_score = %s,
            risk_level = %s,
            header_valid = %s,

            dkim_check_valid = %s,
            dkim_check_details = %s,

            spf_check_valid = %s,
            spf_check_details = %s,

            domain_alignment_check_valid = %s,
            domain_alignment_check_details = %s,

            reply_to_check_valid = %s,
            reply_to_check_details = %s,

            header_warnings = %s
        WHERE message_id = %s
    """

    db.cursor.execute(query, (
        result.get("prediction"),
        result.get("winner_model"),
        result.get("winner_probability"),

        result.get("ensemble_predicted_label"),
        result.get("ensemble_prob_class0"),
        result.get("ensemble_prob_class1"),

        result.get("nb_pred"),
        result.get("rf_pred"),
        result.get("xgb_pred"),
        result.get("knn_pred"),
        result.get("logreg_pred"),

        result.get("nb_prob_class0"),
        result.get("rf_prob_class0"),
        result.get("xgb_prob_class0"),
        result.get("knn_prob_class0"),
        result.get("logreg_prob_class0"),

        result.get("nb_prob_class1"),
        result.get("rf_prob_class1"),
        result.get("xgb_prob_class1"),
        result.get("knn_prob_class1"),
        result.get("logreg_prob_class1"),

        json.dumps(result.get("top_5_words_from_nb", {})),
        json.dumps(result.get("top_5_words_from_rf", {})),
        json.dumps(result.get("top_5_words_from_xgb", {})),
        json.dumps(result.get("top_5_words_from_knn", {})),
        json.dumps(result.get("top_5_words_from_logreg", {})),

        result.get("risk_score"),
        result.get("risk_level"),
        result.get("valid"),

        result.get("checks", {}).get("dkim", {}).get("valid"),
        result.get("checks", {}).get("dkim", {}).get("details"),

        result.get("checks", {}).get("spf", {}).get("valid"),
        result.get("checks", {}).get("spf", {}).get("details"),

        result.get("checks", {}).get("domain_alignment", {}).get("valid"),
        result.get("checks", {}).get("domain_alignment", {}).get("details"),

        result.get("checks", {}).get("reply_to_risk", {}).get("valid"),
        result.get("checks", {}).get("reply_to_risk", {}).get("details"),

        json.dumps(result.get("warnings", [])),

        result.get("message_id")
    ))

    db.conn.commit()

def run_ml_model(email):
     return predict_phishing(
        id=email.get("id", ""),
        subject=email.get("subject", ""),
        body=email.get("body", "")
    )


def validate_headers(email):
    return validate_email_headers(
        dkim_signature=email.get("dkim", ""),
        received_spf=email.get("spf", ""),
        from_header=email.get("sender", ""),
        return_path=email.get("return_path", ""),
        reply_to="",
        message_id=email.get("message_id", "N/A")
    )

