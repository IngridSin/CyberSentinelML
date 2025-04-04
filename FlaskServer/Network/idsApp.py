from FlaskServer.Common.db import DB
from dotenv import load_dotenv
import os


load_dotenv()

if __name__ == "__main__":
    db = DB(os.getenv("DB_NAME"))

    try:
        db.test_query()
        db.listen_for_new_packets()
        # db.close()
    except KeyboardInterrupt:
        db.close()
