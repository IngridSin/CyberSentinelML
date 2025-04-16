import os
import re
import json
import joblib
import nltk
import numpy as np
from nltk.corpus import stopwords
from scipy.sparse import hstack

# === Load Stopwords ===
nltk.download('stopwords')
stop_words = set(stopwords.words('english'))

# === Clean Input Text ===
def clean_text(text):
    if not isinstance(text, str):
        text = str(text)
    text = re.sub(r'\W', ' ', text)
    return ' '.join(word for word in text.split() if word.lower() not in stop_words)

# === Load Model Pipeline ===
#pipeline_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '../../IM/scripts/phishing_pipeline.joblib'))
pipeline_path = os.path.abspath(os.path.join(os.path.dirname(__file__), 'phishing_pipeline.joblib'))

if not os.path.exists(pipeline_path):
    raise FileNotFoundError(f"Pipeline not found at {pipeline_path}")

pipeline = joblib.load(pipeline_path)
ensemble_model = pipeline['ensemble_model']
tfidf_subject = pipeline['tfidf_subject']
tfidf_body = pipeline['tfidf_body']
models = ensemble_model.named_estimators_

# === Predict Function ===
def predict_phishing(id, subject, body):
    # Clean input
    cleaned_subject = clean_text(subject)
    cleaned_body = clean_text(body)

    # Vectorize using loaded TF-IDF vectorizers
    X_subject = tfidf_subject.transform([cleaned_subject])
    X_body = tfidf_body.transform([cleaned_body])
    X_combined = hstack((X_subject, X_body))

    # Ensemble prediction
    prediction = ensemble_model.predict(X_combined)[0]
    ensemble_probs = ensemble_model.predict_proba(X_combined)[0]
    ensemble_prob_class0 = ensemble_probs[0]
    ensemble_prob_class1 = ensemble_probs[1]

    # Individual model predictions and probabilities
    model_scores = {}
    for name, model in models.items():
        pred = model.predict(X_combined)[0]
        probs = model.predict_proba(X_combined)[0]
        prob_class_0 = probs[0]
        prob_class_1 = probs[1]

        model_scores[f"{name.lower()}_pred"] = pred
        model_scores[f"{name.lower()}_prob_class0"] = prob_class_0
        model_scores[f"{name.lower()}_prob_class1"] = prob_class_1

    # Determine winner
    winner_key = max(
        model_scores,
        key=lambda k: model_scores[f"{k.split('_')[0]}_prob_class1"] if model_scores[k] == prediction else -1
    )
    winner = winner_key.split('_')[0]
    winner_probability = (
        model_scores[f"{winner}_prob_class1"] if prediction == 1
        else model_scores[f"{winner}_prob_class0"]
    )

    # Get feature names
    subject_features = tfidf_subject.get_feature_names_out()
    body_features = tfidf_body.get_feature_names_out()
    all_features = np.concatenate([subject_features, body_features])

    # Get feature contributions for all models
    X_combined_dense = X_combined.toarray()[0]
    all_model_contributions = {}

    for name, model in models.items():
        feature_contributions = {}

        if hasattr(model, 'feature_importances_'):  # Random Forest
            importances = model.feature_importances_
            for i, (feature, importance) in enumerate(zip(all_features, importances)):
                if X_combined_dense[i] > 0:
                    feature_contributions[feature] = importance * X_combined_dense[i]
        elif hasattr(model, 'coef_'):  # Linear SVM
            coef = model.coef_[0] if len(model.coef_.shape) > 1 else model.coef_
            for i, (feature, coef_val) in enumerate(zip(all_features, coef)):
                if X_combined_dense[i] > 0:
                    feature_contributions[feature] = coef_val * X_combined_dense[i]
        elif hasattr(model, 'feature_log_prob_'):  # Naive Bayes
            log_probs = model.feature_log_prob_[1]  # Phishing class
            for i, (feature, log_prob) in enumerate(zip(all_features, log_probs)):
                if X_combined_dense[i] > 0:
                    feature_contributions[feature] = log_prob * X_combined_dense[i]

        # Get top 5 features for this model
        top_features = sorted(feature_contributions.items(), key=lambda x: abs(x[1]), reverse=True)[:5]
        all_model_contributions[f"top_5_words_from_{name.lower()}"] = {feature: float(score) for feature, score in top_features}

    # Build result dictionary
    result = {
        "id": id,
        "subject": subject,
        "body": body,
        "prediction": int(prediction),
        "winner_model": winner,
        "winner_probability": float(winner_probability),
        "ensemble_predicted_label": int(prediction),
        "ensemble_prob_class0": float(ensemble_prob_class0),
        "ensemble_prob_class1": float(ensemble_prob_class1),
        **{f"{name.lower()}_pred": int(model_scores[f"{name.lower()}_pred"])
           for name in models.keys()},
        **{f"{name.lower()}_prob_class0": float(model_scores[f"{name.lower()}_prob_class0"])
           for name in models.keys()},
        **{f"{name.lower()}_prob_class1": float(model_scores[f"{name.lower()}_prob_class1"])
           for name in models.keys()},
        **all_model_contributions  # Add top 5 words from all models
    }

    return result

