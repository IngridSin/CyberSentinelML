import pandas as pd
import numpy as np
import joblib
import os

def predict_packets(df ):

    """
    Predicts attack/benign labels on a given CIC-IDS2017 dataframe using a saved XGBoost pipeline.

    Args:
        df (pd.DataFrame): Raw dataframe (e.g., loaded from a CSV).
        model_path (str): Path to the saved joblib model pipeline.

    Returns:
        pd.DataFrame: Results including Packet_ID, prediction, Flow ID, Timestamp, and optionally True Label.
    """
    model_path = os.path.join(os.path.dirname(__file__), "xgboost_pipeline.joblib")

    # Step 1: Clean and prepare
    df = df.copy()
    df.columns = df.columns.str.strip()
    df = df.dropna(how='all').reset_index(drop=True)

    if 'Timestamp' in df.columns and not pd.api.types.is_datetime64_any_dtype(df['Timestamp']):
        df['Timestamp'] = pd.to_datetime(df['Timestamp'], dayfirst=True, errors='coerce')

    if 'Label' in df.columns:
        df['Label'] = df['Label'].astype(str).str.strip().str.upper()
        df['Label'] = df['Label'].apply(lambda x: 1 if x != 'BENIGN' else 0)

    # Step 2: Prepare features
    drop_cols = ['Source IP', 'Destination IP', 'Init_Win_bytes_forward',
                 'Init_Win_bytes_backward', 'Flow ID', 'Timestamp', 'Label', 'Packet_ID']
    X = df.drop(columns=drop_cols, errors='ignore')

    # Sanitize input (inf, nan, overflows)
    X.replace([np.inf, -np.inf], np.nan, inplace=True)
    X.fillna(X.mean(numeric_only=True), inplace=True)
    X = X.clip(lower=-1e308, upper=1e308)

    # Step 3: Load model
    pipeline = joblib.load(model_path)
    expected_cols = pipeline.named_steps['scaler'].feature_names_in_

    # Ensure feature alignment
    X = X[expected_cols]

    # Step 4: Predict
    preds = pipeline.predict(X)

    # Step 5: Add predictions directly to the original df
    df['XGBoost Prediction'] = preds

    return df

