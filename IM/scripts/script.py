import pandas as pd
import re
import nltk
from nltk.corpus import stopwords
from sklearn.model_selection import train_test_split
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.naive_bayes import MultinomialNB
from sklearn.metrics import accuracy_score, classification_report, confusion_matrix, roc_curve, auc
import matplotlib.pyplot as plt
import seaborn as sns

# Load the dataset
df_2 = pd.read_csv("../Data/data.csv")

# Text Preprocessing Function
nltk.download('stopwords') # Words like "the", "is", "in", "and", "to"
stop_words = set(stopwords.words('english'))

# Remove empty Email Text row
df_2 = df_2[df_2['Email Text'].str.strip() != ""] # Remove both empty strings and strings with only spaces
df_2 = df_2.dropna(subset=['Email Text']) # Remove rows that is NaN

def clean_text(text):
    text = text.lower()  # Convert to lowercase
    text = re.sub(r'\W', ' ', text)  # Remove special characters
    text = re.sub(r'\s+', ' ', text)  # Remove extra spaces
    text = ' '.join(word for word in text.split() if word not in stop_words)  # Remove stopwords
    return text

df_2['cleaned_text'] = df_2['Email Text'].apply(clean_text)

# Convert Email Type to boolean
df_2['Email Type'] = df_2['Email Type'].map({'Safe Email': 0, 'Phishing Email': 1})

df_2 = df_2.dropna(subset=['Email Type']) # Remove empty rows

# Split dataset
X_train, X_test, y_train, y_test = train_test_split(df_2['cleaned_text'], df_2['Email Type'], test_size=0.2, random_state=42)

# Convert text to numerical features using TF-IDF
vectorizer = TfidfVectorizer() # Convert text to a matrix of TF-IDF features
X_train_tfidf = vectorizer.fit_transform(X_train) # Learn the vocabulary and the weights of words, Convert into a numerical matrix
X_test_tfidf = vectorizer.transform() # Convert into a numerical matrix

# Train the model
model = MultinomialNB() # For text classification
model.fit(X_train_tfidf, y_train)

# Predictions
y_pred = model.predict(X_test_tfidf)

# Evaluate the model
print("Accuracy:", accuracy_score(y_test, y_pred))
print("Classification Report:\n", classification_report(y_test, y_pred))

# Generate confusion matrix
cm = confusion_matrix(y_test, y_pred)
plt.figure(figsize=(6,4))
sns.heatmap(cm, annot=True, fmt='d', cmap='Blues', xticklabels=['Safe Email', 'Phishing Email'], yticklabels=['Safe Email', 'Phishing Email'])
plt.xlabel('Predicted')
plt.ylabel('Actual')
plt.title('Confusion Matrix')
plt.show()

# Generate ROC Curve
fpr, tpr, _ = roc_curve(y_test, model.predict_proba(X_test_tfidf)[:, 1])
roc_auc = auc(fpr, tpr)
plt.figure(figsize=(6,4))
plt.plot(fpr, tpr, color='blue', label=f'ROC curve (area = {roc_auc:.2f})')
plt.plot([0, 1], [0, 1], color='gray', linestyle='--')
plt.xlabel('False Positive Rate')
plt.ylabel('True Positive Rate')
plt.title('Receiver Operating Characteristic (ROC) Curve')
plt.legend()
plt.show()