FROM python:3.13

WORKDIR /app
COPY . .
COPY Emails/phishing_pipeline.joblib Emails/phishing_pipeline.joblib
COPY Network/xgboost_pipeline.joblib Network/xgboost_pipeline.joblib


RUN pip install --upgrade pip
RUN pip install -r requirements.txt

ENV PYTHONPATH=/app

EXPOSE 5000

CMD ["python", "app.py"]
